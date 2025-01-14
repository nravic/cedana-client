package cmd

// This file contains all the daemon-related commands when starting `cedana daemon ...`

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	task "buf.build/gen/go/cedana/task/protocolbuffers/go"
	"cloud.google.com/go/pubsub"
	"github.com/cedana/cedana/pkg/api"
	"github.com/cedana/cedana/pkg/api/services"
	"github.com/cedana/cedana/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start daemon for cedana client. Must be run as root, needed for all other cedana functionality.",
}

var (
	DEFAULT_PORT      uint32 = 8080
	ASR_POLL_INTERVAL        = 60 * time.Second
)

var startDaemonCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the rpc server. To run as a daemon, use the provided script (systemd) or use systemd/sysv/upstart.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if os.Getuid() != 0 {
			return fmt.Errorf("daemon must be run as root")
		}

		if viper.GetBool("profiling_enabled") {
			go startProfiler()
		}

		var err error

		gpuEnabled, _ := cmd.Flags().GetBool(gpuEnabledFlag)
		vsockEnabled, _ := cmd.Flags().GetBool(vsockEnabledFlag)
		remotingEnabled, _ := cmd.Flags().GetBool(remotingEnabledFlag)
		port, _ := cmd.Flags().GetUint32(portFlag)
		metricsEnabled, _ := cmd.Flags().GetBool(metricsEnabledFlag)
		jobServiceEnabled, _ := cmd.Flags().GetBool(jobServiceFlag)

		if remotingEnabled {
			ctx, err = awsSetup("", ctx, false)
			if err != nil {
				return err
			}
		}

		cedanaURL := viper.GetString("connection.cedana_url")
		if cedanaURL == "" {
			log.Warn().Msg("CEDANA_URL or CEDANA_AUTH_TOKEN unset, certain features may not work as expected.")
			cedanaURL = "unset"
		}

		log.Info().Str("version", rootCmd.Version).Msg("starting daemon")

		// poll for otel signoz logging
		otel_enabled := viper.GetBool("otel_enabled")
		if otel_enabled {
			_, err := utils.InitOtel(ctx, rootCmd.Version)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to initialize otel")
				// fallback to noop
				log.Warn().Msg("Falling back to noop tracer")
				utils.InitOtelNoop()
			}
		} else {
			utils.InitOtelNoop()
		}
		if metricsEnabled {
			pollForAsrMetricsReporting(ctx, port)
		}

		err = api.StartServer(ctx, &api.ServeOpts{
			GPUEnabled:        gpuEnabled,
			VSOCKEnabled:      vsockEnabled,
			CedanaURL:         cedanaURL,
			MetricsEnabled:    metricsEnabled,
			JobServiceEnabled: jobServiceEnabled,
			Port:              port,
		})
		if err != nil {
			log.Error().Err(err).Msgf("stopping daemon")
			return err
		}

		return nil
	},
}

func getenv(k, d string) string {
	if s, f := os.LookupEnv(k); f {
		return s
	}
	return d
}

type AWSCredentials struct {
	AWS_ACCESS_KEY_ID     string `json:"AWS_ACCESS_KEY_ID"`
	AWS_DEFAULT_REGION    string `json:"AWS_DEFAULT_REGION"`
	AWS_SECRET_ACCESS_KEY string `json:"AWS_SECRET_ACCESS_KEY"`
}

func awsCredentialsSetup() error {
	cedana_auth_token, ok := os.LookupEnv("CEDANA_AUTH_TOKEN")
	if !ok {
		return fmt.Errorf("CEDANA_AUTH_TOKEN not set")
	}
	cedana_url, ok := os.LookupEnv("CEDANA_URL")
	if !ok {
		return fmt.Errorf("CEDANA_URL not set")
	}

	// construct + send request
	url := fmt.Sprintf("%s/streaming/aws/credentials", cedana_url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+cedana_auth_token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting AWS credentials: %d", resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// unmarshal response
	var creds AWSCredentials
	err = json.Unmarshal([]byte(respBody), &creds)
	if err != nil {
		log.Err(err).Msg("Error unmarshaling JSON")
		return err
	}

	// set env vars
	err = os.Setenv("AWS_ACCESS_KEY_ID", creds.AWS_ACCESS_KEY_ID)
	if err != nil {
		log.Err(err).Msg("Error setting AWS_ACCESS_KEY_ID")
		return err
	}
	err = os.Setenv("AWS_DEFAULT_REGION", creds.AWS_DEFAULT_REGION)
	if err != nil {
		log.Err(err).Msg("Error setting AWS_DEFAULT_REGION")
		return err
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", creds.AWS_SECRET_ACCESS_KEY)
	if err != nil {
		log.Err(err).Msg("Error setting AWS_SECRET_ACCESS_KEY")
		return err
	}
	return nil
}

func awsSetup(bucket string, ctx context.Context, clear bool) (context.Context, error) {
	env_set := os.Getenv("AWS_DEFAULT_REGION") != "" && os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != ""
	file_set := os.Getenv("AWS_CONFIG_FILE") != "" && os.Getenv("AWS_SHARED_CREDENTIALS_FILE") != ""
	if !env_set && !file_set {
		err := awsCredentialsSetup()
		if err != nil {
			log.Error().Err(err).Msg("Failed to setup AWS credentials")
			return ctx, fmt.Errorf("Failed to setup AWS credentials")
		}
	}
	if os.Getenv("AWS_DEFAULT_REGION") == "" && os.Getenv("AWS_CONFIG_FILE") == "" {
		return ctx, fmt.Errorf("Please set environment variable AWS_DEFAULT_REGION, or set AWS_CONFIG_FILE to absolute path of AWS config file (~/.aws/config).")
	}
	if !((os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_SECRET_ACCESS_KEY") == "") || os.Getenv("AWS_SHARED_CREDENTIALS_FILE") == "") {
		return ctx, fmt.Errorf("Please set environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY, or set AWS_SHARED_CREDENTIALS_FILE to absolute path of AWS credentials file (~/.aws/credentials).")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return ctx, fmt.Errorf("Unable to load AWS configuration")
	}
	if cfg.Region == "" {
		return ctx, fmt.Errorf("AWS region not configured properly, please specify in AWS config file (~/.aws/config) or environment variable AWS_DEFAULT_REGION.")
	}
	_, err = cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return ctx, fmt.Errorf("Failed to load AWS credentials, please specify in AWS credentials file (~/.aws/credentials) or environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.")
	}

	s3Client := s3.NewFromConfig(cfg)
	if bucket != "" {
		_, err = s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return ctx, err
		}
	}
	if clear { // clear bucket on dump
		var objectsToDelete []types.ObjectIdentifier
		var page *s3.ListObjectsV2Output
		paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
		for paginator.HasMorePages() {
			page, err = paginator.NextPage(ctx)
			if err != nil {
				return ctx, fmt.Errorf("failed to list objects in bucket: %v", err)
			}
		}
		for _, obj := range page.Contents {
			objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}
		for i := 0; i < len(objectsToDelete); i += 1000 {
			end := i + 1000
			if end > len(objectsToDelete) {
				end = len(objectsToDelete)
			}
			_, err = s3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &types.Delete{
					Objects: objectsToDelete[i:end],
				},
			})
			if err != nil {
				return ctx, fmt.Errorf("failed to delete objects: %v", err)
			}
		}
		output, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return ctx, fmt.Errorf("failed to list objects in bucket: %v", err)
		}
		if len(output.Contents) != 0 {
			return ctx, fmt.Errorf("failed to clear bucket: %v objects remain", len(output.Contents))
		}
	}
	return ctx, nil
}

func gcloudAdcSetup(ctx context.Context) error {
	adcPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if adcPath == "" {
		// set env if not present
		// default to root /gcloud-credentials.json
		adcPath = "/gcloud-credentials.json"
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", adcPath)
	}
	if _, err := os.Stat(adcPath); err == nil {
		// already present skip
		return nil
	}
	cedanaURL := viper.GetString("connection.cedana_url")
	url := cedanaURL + "/k8s/gcloud/serviceaccount"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("connection.cedana_auth_token")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(adcPath, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func pollForAsrMetricsReporting(ctx context.Context, port uint32) {
	// polling for ASR
	go func() {
		// setup GCLOUD_JSON
		err := gcloudAdcSetup(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to setup gcloud ADC, disabling reporting")
			return
		}
		// end
		log.Info().Msg("start pushing asr metrics")
		client, err := pubsub.NewClient(ctx, getenv("GOOGLE_CLOUD_PROJECT", "prod-data-438318"))
		if err != nil {
			log.Error().Msgf("Failed to create Pub/Sub client: %v", err)
			return
		}
		defer client.Close()
		manager, err := api.SetupCadvisor(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to setup cadvisor")
			return
		}

		macAddr, _ := utils.GetMACAddress()
		hostname, _ := os.Hostname()
		v, _ := mem.VirtualMemory()
		pmem := fmt.Sprintf("%d", v.Total/(1024*1024*1024)) // in GB
		url := viper.GetString("connection.cedana_url")

		topic := client.Topic("asr-metrics")
		time.Sleep(10 * time.Second)
		for {
			conts, err := api.GetContainerInfo(ctx, manager)
			if err != nil {
				log.Error().Msgf("error getting info: %v", err)
				return
			}
			b, err := json.Marshal(conts)
			// Publish a message
			result := topic.Publish(ctx, &pubsub.Message{
				Data: b,
				Attributes: map[string]string{
					"mac":      macAddr,
					"hostname": hostname,
					"mem":      pmem,
					"url":      url,
				},
			})
			// Get the server-assigned message ID
			id, err := result.Get(ctx)
			if err != nil {
				log.Error().Msgf("Failed to publish message: %v", err)
			}
			log.Info().Msgf("Published message with ID: %v\n", id)
			time.Sleep(ASR_POLL_INTERVAL)
		}
	}()
}

var checkDaemonCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if daemon is running and healthy",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetUint32(portFlag)

		cts, err := services.NewClient(port)
		if err != nil {
			log.Error().Err(err).Msg("error creating client")
			return err
		}

		defer cts.Close()

		// regular health check
		healthy, err := cts.HealthCheck(cmd.Context())
		if err != nil {
			return err
		}
		if !healthy {
			return fmt.Errorf("health check failed")
		}

		// Detailed health check. Need to grab uid and gid to start
		// controller properly and with the right perms.
		var uid int32
		var gid int32
		var groups []int32 = []int32{}

		uid = int32(os.Getuid())
		gid = int32(os.Getgid())
		groups_int, err := os.Getgroups()
		if err != nil {
			return fmt.Errorf("error getting user groups: %v", err)
		}
		for _, g := range groups_int {
			groups = append(groups, int32(g))
		}

		req := &task.DetailedHealthCheckRequest{
			UID:    uid,
			GID:    gid,
			Groups: groups,
		}

		resp, err := cts.DetailedHealthCheck(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("health check failed: %v", err)
		}

		if len(resp.UnhealthyReasons) > 0 {
			return fmt.Errorf("health failed with reasons: %v", resp.UnhealthyReasons)
		}

		fmt.Println("All good.")
		fmt.Println("Cedana version: ", rootCmd.Version)
		fmt.Println("CRIU version: ", resp.HealthCheckStats.CriuVersion)
		if resp.HealthCheckStats.GPUHealthCheck != nil {
			prettyJson, err := json.MarshalIndent(resp.HealthCheckStats.GPUHealthCheck, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println("GPU support: ", string(prettyJson))
		}

		return nil
	},
}

// Used for debugging and profiling only!
func startProfiler() {
	utils.StartPprofServer()
}

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(startDaemonCmd)
	daemonCmd.AddCommand(checkDaemonCmd)
	startDaemonCmd.Flags().BoolP(gpuEnabledFlag, "g", false, "start daemon with GPU support")
	startDaemonCmd.Flags().Bool(vsockEnabledFlag, false, "start daemon with vsock support")
	startDaemonCmd.Flags().BoolP(remotingEnabledFlag, "r", false, "start daemon with direct remoting support")
	startDaemonCmd.Flags().BoolP(metricsEnabledFlag, "m", false, "enable metrics")
	startDaemonCmd.Flags().Bool(jobServiceFlag, false, "enable job service")
}
