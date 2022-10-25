package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var kube_url = "https://localhost:38045"

func main() {
	assemblingArt()
}

// gets the amount of nodes in the cluster
func getNodeCount() (nodeCount int) {

	//loads certificate
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	//initializes client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	//the apiEndpoint gets declared as a variable to ease logging
	apiEndpoint := "/api/v1/nodes"

	//fetches data about nodes
	resp, err := client.Get(kube_url + apiEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//continues if HTTP status code is ok, logs an error otherwise
	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the node count
		nodeCount = strings.Count(bodyString, "managedFields")

	} else {
		//logs the error
		log.Println("Connection error:", kube_url+apiEndpoint)
	}

	//returns the node count
	//the functions output gets used somewhere else
	return nodeCount
}

// gets the amount of nodes in the cluster
func getNamespaceCount() (namespaceCount int) {

	//loads the certs
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	//initializes the client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	//calls the api endpoint for namespaces
	apiEndpoint := "/api/v1/namespaces"
	resp, err := client.Get(kube_url + apiEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//continues if the HTTP status code is ok
	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per namespace in the json file you get when calling the url, so it is getting used to get the namespace count
		namespaceCount = strings.Count(bodyString, "managedFields")

	} else {
		//when connection to API fails,logs error containing the full endpoint URL
		log.Println("Connection error:", kube_url+apiEndpoint)
	}

	return namespaceCount
}

func getPodCount() (podCount int) {

	//loads certificates needed for HTTP client
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	//initializes HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	//calls the API endpoint for pods
	apiEndpoint := "/api/v1/pods"
	resp, err := client.Get(kube_url + apiEndpoint)

	//errors out if that fails
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//continues if the HTTP status code is ok
	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//managedFields only appears once per node in the json file you get when calling the url, so it is getting used to get the pod count
		podCount = strings.Count(bodyString, "managedFields")

	} else {
		//logs an error containing the full API URL if the connection fails
		log.Println("connection error:", kube_url+apiEndpoint)
	}

	return podCount
}

func getDistro() (majorVersion string, minorVersion string, distro string) {

	//loads the certificates required
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	//initializes the HTTP client with the certificates provided
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	apiEndpoint := "/version"
	resp, err := client.Get(kube_url + apiEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//continues if the HTTP status code is ok
	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(bodyBytes)

		//parses the field "gitVersion" with the git version the kubernetes distribution uses
		//this can be used to detect which distro you are running
		//this needs to get converted into a string again, thats why theres a .String() at the end
		gitVersion := gjson.Get(bodyString, "gitVersion").String()

		//parses the field "major" containing the major Kubernetes Version
		majorVersion = gjson.Get(bodyString, "major").String()
		minorVersion = gjson.Get(bodyString, "minor").String()

		if strings.Contains(gitVersion, "k3s") {
			//currently only detection for k3s is implemented
			distro = "k3s"
		} else {
			//it assumes that every Kubernetes distro that doesn't have k3s in their git version is normal Kubernetes
			distro = "k8s"
		}

	} else {
		//errors out when the /version endpoint can't be reached
		log.Println("Connection error:", kube_url+apiEndpoint)
	}
	return majorVersion, minorVersion, distro

}

func getUsedIngress() (ingressUsed string) {

	//loads the client certificates
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatal(err)
	}

	//initializes the client with the client certificates from above
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	//the detection methods for nginx and traefik are different because Traefik ingresses aren't defined with the normal networking api, they use a CRD

	//to detect a traefik install, all CRDs get fetched
	respCouldContainTraefik, traefikErr := client.Get(kube_url + "/apis/apiextensions.k8s.io/v1/customresourcedefinitions")

	//to detect a Nginx-Ingress install, all ingressclasses get fetched
	respCouldContainNginx, nginxErr := client.Get(kube_url + "/apis/networking.k8s.io/v1/ingressclasses")

	//errors out if any of the requests fail
	if traefikErr != nil {
		log.Fatal(err)
	}

	if nginxErr != nil {
		log.Fatal(err)
	}

	defer respCouldContainTraefik.Body.Close()
	defer respCouldContainNginx.Body.Close()

	//only continues if both request return a HTTP status code 200
	if respCouldContainTraefik.StatusCode == http.StatusOK && respCouldContainNginx.StatusCode == http.StatusOK {

		couldContainTraefik, errTraefik := io.ReadAll(respCouldContainTraefik.Body)
		couldContainNginx, errNginx := io.ReadAll(respCouldContainNginx.Body)

		if errTraefik != nil {
			log.Fatal(err)
		}
		if errNginx != nil {
			log.Fatal(err)
		}

		//kubefetch can show if you use multiple ingresses, so the string for the ingresses can be composed out of multiple values
		//the bool first-append is neccesary because if gets used to detect if a comma should be put at the front of the string ingressUsed
		var firstAppend bool

		firstAppend = true

		if strings.Contains(string(couldContainTraefik), "traefik.containo.us") {
			ingressUsed = "Traefik"
			firstAppend = false
		} else {
		}

		if strings.Contains(string(couldContainNginx), "nginx") {
			//checks if the string ingressUsed has already been appended to

			//if that happens, it will add a comma, otheriwse it won't add one
			if firstAppend == false {
				ingressUsed = ingressUsed + ", "
			} else {
				firstAppend = true
			}
			ingressUsed = ingressUsed + "Nginx"
		} else {
		}

	}
	return ingressUsed

}

func assemblingArt() {

	//initializes all variables
	var majorVersion string
	var minorVersion string
	var distro string

	//this initializes the variables that get used for printing the info
	var distroArt [17]string
	var asciiArtColor string

	if distro == "k8s" {

		distroArt[0] = "		0MMMNNMMMO               	"
		distroArt[1] = "          .MMMWKko:::cokKWMMM.          	"
		distroArt[2] = "      dMMMN0xl:::::cc:::::lx0NMMMo      	"
		distroArt[3] = "    0MXkdc:::::::::K0:::::::::cdOXMO    	"
		distroArt[4] = "   'MO::::::::::codXXdoc::::::::::0M.   	"
		distroArt[5] = "   OWc:::oXkclONWX0MM0XWXklckXl:::cMk   	"
		distroArt[6] = "  .MO::::::dMMWx::cMM:::kMMMd::::::0M.  	"
		distroArt[7] = "  OMc::::::XMkKMNkkMMkONM0kMX::::::lMx  	"
		distroArt[8] = " .MO::::::xMO::cKMM00MMKc::0Md::::::0M. 	"
		distroArt[9] = " kMc::::::OMX0XWMMWdxWMMWX0XMk::::::lMx 	"
		distroArt[10] = " M0::::kKOOWMkolcXMWWMXclokMWkOKk::::KM 	"
		distroArt[11] = " .No:::::::lNWxcOMKccXMkckWNl:::::::oN. 	"
		distroArt[12] = "   c0c:::::::xXWMWxddxWMWXx:::::::c0:   	"
		distroArt[13] = "     Ox:::::::lWkkO00OkkWl:::::::kk     	"
		distroArt[14] = "       Xo::::cKd::::::::xK:::::oN       	"
		distroArt[15] = "        :0c::::::::::::::::::lK,        	"
		distroArt[16] = "          kklcccccccccccccclOx          	"
		asciiArtColor = "34"

	} else if distro == "k3s" {

		distroArt[0] = " .kWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWk.  	"
		distroArt[1] = ":WMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWc	"
		distroArt[2] = "WMMMMMMMMMMMMMMMMM    MMMMMMMMMMMMMMMMMM	"
		distroArt[3] = "WMMMMMMMMMMMMMMMMX    MMMMMMMMMMMMMMMMMM	"
		distroArt[4] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[5] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[6] = "WMMMMMMMMMMMMMMMMK    MMMMMMMMMMMMMMMMMM	"
		distroArt[7] = "WMMMMMMMMMMMMMMMMW    MMMMMMMMMMMMMMMMMM	"
		distroArt[8] = "WMMMMMMMMMMMMMMMMMW00WMMMMMMMMMMMMMMMMMM	"
		distroArt[9] = "WMMMMMMMMMMMMM   lMMMMc   MMMMMMMMMMMMMM	"
		distroArt[10] = "WMMMMMMMMM.       MMMM       .MMMMMMMMMM	"
		distroArt[11] = "WMMMMMx        'lXMMMMKl.        kMMMMMM	"
		distroArt[12] = "WMMMMM     .cOWMMMMMMMMMMWkc.    .MMMMMM	"
		distroArt[13] = "XMMMMMXlcxNMMMMMMMMMMMMMMMMMMNxclXMMMMMN	"
		distroArt[14] = " oMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMo 	"
		distroArt[15] = "    ''''''''''''''''''''''''''''''''    	"
		distroArt[16] = "										 	"
		asciiArtColor = "93"

	} else {

		distroArt[0] = "		0MMMNNMMMO               	"
		distroArt[1] = "          .MMMWKko:::cokKWMMM.          	"
		distroArt[2] = "      dMMMN0xl:::::cc:::::lx0NMMMo      	"
		distroArt[3] = "    0MXkdc:::::::::K0:::::::::cdOXMO    	"
		distroArt[4] = "   'MO::::::::::codXXdoc::::::::::0M.   	"
		distroArt[5] = "   OWc:::oXkclONWX0MM0XWXklckXl:::cMk   	"
		distroArt[6] = "  .MO::::::dMMWx::cMM:::kMMMd::::::0M.  	"
		distroArt[7] = "  OMc::::::XMkKMNkkMMkONM0kMX::::::lMx  	"
		distroArt[8] = " .MO::::::xMO::cKMM00MMKc::0Md::::::0M. 	"
		distroArt[9] = " kMc::::::OMX0XWMMWdxWMMWX0XMk::::::lMx 	"
		distroArt[10] = " M0::::kKOOWMkolcXMWWMXclokMWkOKk::::KM 	"
		distroArt[11] = " .No:::::::lNWxcOMKccXMkckWNl:::::::oN. 	"
		distroArt[12] = "   c0c:::::::xXWMWxddxWMWXx:::::::c0:   	"
		distroArt[13] = "     Ox:::::::lWkkO00OkkWl:::::::kk     	"
		distroArt[14] = "       Xo::::cKd::::::::xK:::::oN       	"
		distroArt[15] = "        :0c::::::::::::::::::lK,        	"
		distroArt[16] = "          kklcccccccccccccclOx          	"
		asciiArtColor = "34"
	}

	var nodeCount = getNodeCount()
	var podCount = getPodCount()
	var namespaceCount = getNamespaceCount()
	var ingressUsed = getUsedIngress()
	majorVersion, minorVersion, distro = getDistro()
	//fetch all values from the various functions above

	print("\033[" + asciiArtColor + ";1m" + distroArt[0] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mDistro: \033[0m")
	print(distro)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[1] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mVersion: \033[0m")
	print(majorVersion, ".", minorVersion)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[2] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mNode Count: \033[0m")
	print(nodeCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[3] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mPod Count: \033[0m")
	print(podCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[4] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mNamespace Count: \033[0m")
	print(namespaceCount)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[5] + "\033[0m")
	print("\033[" + asciiArtColor + ";1mIngress: \033[0m")
	print(ingressUsed)
	print("\n")

	print("\033[" + asciiArtColor + ";1m" + distroArt[6] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[7] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[8] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[9] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[10] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[11] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[12] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[13] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[14] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[15] + "\033[0m")
	print("\n")
	print("\033[" + asciiArtColor + ";1m" + distroArt[16] + "\033[0m")
	print("\n")
	print("")
}
