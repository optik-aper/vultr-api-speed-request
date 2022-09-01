package main

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/optik-aper/vultr-api-speed-request/pkg/vultr"
	"github.com/vultr/govultr/v2"
)

const (
	requestNum int    = 24
	rateLimit  int    = 100
	retryLimit int    = 1
	label      string = "speedyIP"
	ipRegion   string = "tyo"
)

// Execute the main program
func main() {
	var apiKey string
	if apiKey = os.Getenv("VULTR_API_KEY"); apiKey == "" {
		panic("API KEY missing")
	}

	config := vultr.Config{
		APIKey:     apiKey,
		RateLimit:  rateLimit,
		RetryLimit: retryLimit,
	}

	client, err := config.Init()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "up":
		up(client)
	case "down":
		down(client)
	}
}

func up(client *govultr.Client) {

	wg := sync.WaitGroup{}
	errs := make(chan error, 1)
	done := make(chan bool, 1)

	for i := 0; i < requestNum; i++ {

		wg.Add(1)

		go func(num int) {
			if _, err := client.ReservedIP.Create(context.Background(), &govultr.ReservedIPReq{
				Region: "lax",
				IPType: "v4",
				Label:  label + "_" + strconv.Itoa(num),
			}); err != nil {
				errs <- err
			}

			defer wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case err := <-errs:
		if err != nil {
			panic(err.Error())
		}
	}
}

func down(client *govultr.Client) {

	wg := sync.WaitGroup{}
	errs := make(chan error, 1)
	done := make(chan bool, 1)

	ips, _ := resIPList(context.Background(), client)

	for _, ip := range ips {

		wg.Add(1)

		go func(ipID string) {
			if err := resIPDelete(context.Background(), client, ipID); err != nil {
				errs <- err
			}

			defer wg.Done()
		}(ip.ID)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case err := <-errs:
		if err != nil {
			panic(err.Error())
		}
	}
}

func resIPCreate(ctx context.Context, client *govultr.Client, label string) error {
	_, err := client.ReservedIP.Create(ctx, &govultr.ReservedIPReq{
		Region: ipRegion,
		IPType: "v4",
		Label:  label,
	})

	if err != nil {
		return err
	}
	return nil
}

func resIPDelete(ctx context.Context, client *govultr.Client, id string) error {
	if err := client.ReservedIP.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func resIPList(ctx context.Context, client *govultr.Client) ([]govultr.ReservedIP, error) {
	ips, _, err := client.ReservedIP.List(ctx, nil)

	if err != nil {
		return nil, err
	}

	return ips, nil
}
