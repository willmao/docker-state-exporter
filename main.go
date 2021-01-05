package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	//DockerAPIVersion Docker API Version
	DockerAPIVersion = "1.24"
)

var (
	addr             = flag.String("listen-address", ":9901", "The address to listen on for HTTP requests.")
	refreshInterval  = flag.Float64("refresh-interval", 5, "The interval to fetch container state")
	prefixesToSkip   = flag.String("prefix-to-skip", "k8s_", "Skip collecting container state with specified prefix")
	dockerStateGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "docker_container_state",
		Help: "Docker container state"},
		[]string{"container_name", "state", "exit_code"})
)

func skippedContainer(prefixesToSkip []string, containerName string) bool {
	for _, prefix := range prefixesToSkip {
		if strings.HasPrefix(containerName, "/"+prefix) {
			return true
		}
	}
	return false
}

func init() {
	prometheus.MustRegister(dockerStateGauge)
}

func main() {
	flag.Parse()
	lastTime := time.Now()

	sigs := make(chan os.Signal, 1)

	go func() {
		sig := <-sigs
		log.Printf("receive sig: %s", sig)
		os.Exit(0)
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cli, err := client.NewClientWithOpts(client.WithVersion(DockerAPIVersion))
	defer cli.Close()
	if err != nil {
		panic(err)
	}

	_, err = cli.Info(context.Background())
	if err != nil {
		panic(err)
	}

	log.Println("test docker client connectivity: ok")
	prefixes := strings.Split(*prefixesToSkip, ",")
	go func() {
		for {
			currentTime := time.Now()
			if currentTime.Sub(lastTime).Seconds() > *refreshInterval {
				containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
				if err != nil {
					log.Println(err)
				}
				lastTime = time.Now()
				dockerStateGauge.Reset()
				for _, container := range containers {
					if len(container.Names) == 0 {
						continue
					}

					containerName := container.Names[0]
					containerName = strings.TrimPrefix(containerName, "/")
					if skippedContainer(prefixes, containerName) {
						continue
					}

					exitCode := 0
					containerJSON, err := cli.ContainerInspect(context.Background(), container.ID)
					if err != nil {
						log.Printf("failed to get container exit code, error: %s\n", err.Error())
					} else {
						exitCode = containerJSON.ContainerJSONBase.State.ExitCode
					}

					log.Printf("container name: %s, state: %s, exit code: %d\n", containerName, container.State, exitCode)
					dockerStateGauge.WithLabelValues(containerName, container.State, strconv.Itoa(exitCode)).Set(1)
				}
			}

			time.Sleep(time.Second * 1)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
