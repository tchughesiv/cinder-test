package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/noauth"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
)

const (
	oshost = "https://cinder.openstack.svc"
)

var (
	tokenFile = flag.String("token", "/var/run/secrets/kubernetes.io/serviceaccount/token", "A pod's serviceaccount bearer token")
	caFile    = flag.String("ca", "/var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt", "A PEM eoncoded CA's certificate file.")
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func contains(opts []string, option string) bool {
	for _, a := range opts {
		if a == option {
			return true
		}
	}
	return false
}

func main() {
	var option string
	var volID string
	var csize int
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	opts := []string{"create", "delete", "get"}
	if !contains(opts, option) {
		fmt.Printf("\n")
		fmt.Printf("%#v is not a valid request. Must choose one of these:\n", option)
		for _, opt := range opts {
			fmt.Printf(" - %s\n", opt)
		}
		fmt.Printf("\n")
		os.Exit(1)
	}

	token, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	bearer := map[string]string{"Authorization": "Bearer " + string(token)}

	ao, err := openstack.AuthOptionsFromEnv()

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	provider, err := noauth.NewClient(ao)
	checkErr(err)
	client, err := noauth.NewBlockStorageV2(provider, noauth.EndpointOpts{
		CinderEndpoint: fmt.Sprintf("%s/v2", oshost),
		// CinderEndpoint: os.Getenv("CINDER_ENDPOINT"),
	})

	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	if option == "create" {
		if len(os.Args) < 3 {
			fmt.Printf("\n")
			fmt.Printf("Must pass size of volume to create:")
			fmt.Printf("\n\n")
			fmt.Printf("  cinder-test %v <size>", option)
			fmt.Printf("\n\n")
			os.Exit(1)
		}
		csize, err = strconv.Atoi(os.Args[2])
		checkErr(err)
		vopts := volumes.CreateOpts{Size: csize, VolumeType: "iscsi"}
		vol, err := Create(client, vopts, bearer).Extract()
		checkErr(err)
		log.Printf("%v - Created at %v", vol.ID, vol.CreatedAt)
	}

	if option == "get" {
		if len(os.Args) < 3 {
			fmt.Printf("\n")
			fmt.Printf("Must pass volume ID:")
			fmt.Printf("\n\n")
			fmt.Printf("  cinder-test %v <volume ID>", option)
			fmt.Printf("\n\n")
			os.Exit(1)
		}
		volID = os.Args[2]
		gvol, err := Get(client, volID, bearer).Extract()
		checkErr(err)
		log.Printf("%v", gvol)
	}

	if option == "delete" {
		if len(os.Args) < 3 {
			fmt.Printf("\n")
			fmt.Printf("Must pass volume ID:")
			fmt.Printf("\n\n")
			fmt.Printf("  cinder-test %v <volume ID>", option)
			fmt.Printf("\n\n")
			os.Exit(1)
		}
		volID = os.Args[2]
		result := Delete(client, volID, bearer)
		log.Printf("Deleting %v - %v", volID, result)
	}
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("volumes")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

// Create will create a new Volume based on the values in CreateOpts. To extract
// the Volume object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts volumes.CreateOptsBuilder, bearer map[string]string) (r volumes.CreateResult) {
	b, err := opts.ToVolumeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: bearer,
	})
	return
}

// Delete will delete the existing Volume with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string, bearer map[string]string) (r volumes.DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), &gophercloud.RequestOpts{
		MoreHeaders: bearer,
	})
	return
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string, bearer map[string]string) (r volumes.GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: bearer,
	})
	return
}
