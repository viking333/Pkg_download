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
	"crypto/sha256"
)

func main() {
	logfile_create()
	fmt.Println("****************** File Download ******************\n")
	
	//create downoad folder where the new files will be stored. If the folder already exhist. remove old packages.
	//if _, err := os.Stat("downloads"); os.IsNotExist(err){
	_, err := os.Stat("downloads")
	if err != nil {
		fmt.Println("----------- Creating downloads folder -----------")
		os.Mkdir("downloads", 0700)
	}
	//delete leftowers located in "downloads" folder form previous downloads
	RemoveContent("downloads/")
	
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
		//fmt.Printf("The file located in %s was succesfully downloaded. \n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Reading lines of packages file failed with an error: %s\n", err)
	}
	AV_scan()
}

func download(url string) error {
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
	fmt.Printf("__________ Downloading Package__________\nThe %s was succesfully downloaded. \n", filename)
	//calculate file hash
	//log_output := fmt.SprintF("The file was downloaded from %s\n", filemane)
	logfile_write()
	Create_summary(filename)
	return nil
}

	//av scans
func AV_scan(){
	
	//clam av scan
	fmt.Println("\n****************** Anti-Virus Scan ******************\n")
	fmt.Println("\n----------- Clam AV Scan -----------\n")
	clam_cmd := exec.Command("clamscan", "-i", "downloads/")
	clam_av_scan, err := clam_cmd.CombinedOutput()
	if err != nil{
		log.Fatalf("Clam AV scan failed with an error %s\n", err)
	}
	fmt.Printf("The Clam AV scan was completed successfully and the results of it are: \n%s\n", clam_av_scan)
	
	//Sophos Scan
	fmt.Println("\n----------- Sophos Scan -----------\n")
	sav_cmd := exec.Command("savscan", "-f", "-all", "-rec", "-ss", "-archive", "downloads/")
	sav_av_scan, err := sav_cmd.CombinedOutput()
	if err != nil{
		log.Fatalf("Sophos scan failed with an error %s\n", err)
	}
	fmt.Printf("The Sophos scan was completed.The results are: \n%s\n (No results mean the Sophos did not find any issues with the file)", sav_av_scan)

	
}
	// function used to delete existing files in the folder
func RemoveContent(dir string) {
	d, err := os.Open(dir)
	if err != nil {
		log.Fatalf("Old files removal failed with an error %s\n", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Old files removal failed with an error %s\n", err)
	}
	for _, name := range names {
		file := (dir + name)
		err = os.RemoveAll(file)
		if err != nil {
			log.Fatalf("Old file removal filed with an error %s\n", err)
		}
	}
}

	//this function creates a text file with all the data used 
func Create_summary(filename string) {
	
	path := "downloads/" + filename
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Hashing file opening function failed with an error s%\n", err)
	}
	defer f.Close()
	
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatalf("Hashing function failed with an error s%\n", err)
	}

	output := fmt.Sprintf("The SHA256 hash of %s is : %x\n\n", filename,  hasher.Sum(nil))
	logfile_write(output)
}


func logfile_create() {
	//checking if the log file existing and creating one if it doesn't
	_, err := os.Stat("status.log")
	if os.IsNotExist(err) {
		fmt.Println("----------- Creating log file -----------")
		file, err1 := os.Create("status.log")
		if err1 != nil {
			log.Fatalf("Log file creation process failed with an error %s\n", err1)
		defer file.Close()
		}
	}
	//cleaning up existing file
	err2 := os.Truncate("status.log", 0)
	if err2 != nil {
		log.Fatalf("Log file clearing process failed with an error %s\n", err2)
	}

}

func logfile_write(message string) {
	
	//openining file for writing
	file, err := os.OpenFile("status.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Opening log file for editing failed with an error %s", err)
	}
	defer file.Close()

	//writing message to the log file
	_, err = file.WriteString(message)
	if err != nil {
		log.Fatalf("Writing to the log file failed with an error %s\n", err)
	}
	

	// saving changes

	err = file.Sync()
	if err != nil {
		log.Fatalf("Log file save function failed with an error %s\n", err)
	}
	
}
