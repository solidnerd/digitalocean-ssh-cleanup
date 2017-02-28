package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// TokenSource provides an AccessToken for DigitalOcean API
type TokenSource struct {
	AccessToken string
}

// Token will be called via an interface in the oauth client in line 35
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {

	token := flag.String("token", "mytoken", "DigitalOcean API Token")
	keyname := flag.String("keyname", "docker", "A word that is contained in the keyname")
	flag.Parse()
	tokenSource := &TokenSource{
		AccessToken: *token,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	ctx := context.TODO()
	keys, err := sshKeyList(ctx, client)

	if err != nil {
		log.Fatal(err)
	}

	for _, key := range keys {
		if strings.Contains(key.Name, *keyname) {
			log.Println("Delete Key for: ", key.Name)
			client.Keys.DeleteByID(key.ID)
			time.Sleep(10 * time.Second)
		}
	}
}

func sshKeyList(ctx context.Context, client *godo.Client) ([]godo.Key, error) {
	// create a list to hold our droplets
	list := []godo.Key{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		keys, resp, err := client.Keys.List(opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range keys {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}
