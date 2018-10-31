package main

import(
	"fmt"
	"os"
	"os/exec"
	"io"
	"log"
	"bufio"
	"net/http"
	"strings"
)

func main() {

	// open file
	file, err := os.Open("packages.txt")
	if err != nil {
		log.Fatalf("File opening failed with an error: %s\n", err)
	}
	
	//read lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
	
		//download file 
		err := download(scanner.Text())
		if err != nil {
			log.Fatalf("file download failed with an error: %s\n", err)
		}
		fmt.Printf("The file located in %s was succesfully downloaded. \n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Reading lines of packages file failed with an error: %s\n", err)
	}

	fmt.Println("\nThe files will now get scanned by AV solutions\n")
	
	AV_scan()
}

func download(url string) error{

	
	//get data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Getting data of the site failed with an error: %s\n", err)
	}
	defer resp.Body.Close()

	//create out file
	//follow the redirects and get final address, split it by "/" and then pick the last part which should be the name
	finalURL := resp.Request.URL.String()
	parts := strings.Split(finalURL, "/")
	filename := parts[len(parts)-1]
	finalname := "downloads/" + filename

	out, err := os.Create(finalname)
	if err != nil {
		log.Fatalf("File creation failed with an error: %s\n", err)
	}
	
	//write the body to the file 

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Copying downloaded content to the previoustly created file failed with an error: %s\n", err)
	}
	return nil
}

	//av scans
func AV_scan(){
	//clam av scan
	clam_cmd := exec.Command("clamscan", "-i", "downloads/")
	clam_av_scan, err := clam_cmd.CombinedOutput()
	if err != nil{
		log.Fatalf("Clam AV scan failed with an error %s\n", err)
	}
	fmt.Printf("The Clam AV scan was completed successfully and the results of it are: \n%s\n", clam_av_scan)

	sav_cmd := exec.Command("savscan", "-f", "-all", "-rec", "-ss", "-archive", "downloads/")
	sav_av_scan, err := sav_cmd.CombinedOutput()
	if err != nil{
		log.Fatalf("Sophos scan failed with an error %s\n", err)
	}
	fmt.Printf("The Sophos scan was completed successfully and the results of it are: \n%s\n", sav_av_scan)

	
}
