package main
 
import(
  "fmt"
  "os"
  "os/exec"
  "bytes"
  "net/http"
  "io/ioutil"
  "github.com/jessevdk/go-flags"
)

const(
  WP_DOWNLOAD_URL  = "http://wordpress.org/wordpress-%s.tar.gz"
  WP_VERSIONS_FILE = "https://gist.github.com/sosedoff/3730299d7c4ef0c5bc70/raw/wp-versions.txt"
)

func fileExists(path string) bool {
  _, err := os.Stat(path)
 
  if err == nil { return true }
  if os.IsNotExist(err) { return false }
 
  return false
}
 
func run(command string) (string, error) {
  var output bytes.Buffer
 
  cmd := exec.Command("bash", "-c", command)
 
  cmd.Stdout = &output
  cmd.Stderr = &output

  fmt.Println("Running:", command)
 
  err := cmd.Run()
  return output.String(), err
}
 
func wp_core_list() {
  resp, err := http.Get(WP_VERSIONS_FILE)

  if err != nil {
    fmt.Println("Unable to get list of versions...")
    return
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)

  fmt.Println("Available Wordpress versions:")
  fmt.Println(string(body))
}

func wp_core_download(version string, path string) {
  url := fmt.Sprintf(WP_DOWNLOAD_URL, version)
  temp := "/tmp/wordpress.tar.gz"

  if fileExists(temp) {
    run(fmt.Sprintf("rm -f %s", temp))
  }
 
  fmt.Println("Downloading core")
  cmd := fmt.Sprintf("wget -O %s %s", temp, url)
  out, err := run(cmd)

  if err != nil {
    fmt.Println("Unable to download Wordpress Core!")
    fmt.Println(out)
    os.Exit(1)
  }
 
  fmt.Println("Extracting core")
 
  if fileExists(path) {
    fmt.Println("Removing and old directory")
    run(fmt.Sprintf("rm -rf %s", path))
  }
 
  _, err = run(fmt.Sprintf("cd /tmp && tar -zxf %s", temp))
  if err != nil {
    fmt.Println("Failed to extract wordpress core:", err)
    os.Exit(1)
  }
}

func handle_command(command string) {
  if command == "core:list" {
    wp_core_list()
  }

  if command == "core:download" {
    var opts struct {
      Version string `short:"v" long:"version" description:"Core version"`
      Path string    `short:"p" long:"path" description:"Path to install"`
    }

    _, err := flags.ParseArgs(&opts, os.Args)
    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    }

    if len(opts.Version) == 0 {
      opts.Version = "latest"
    }

    if len(opts.Path) == 0 {
      fmt.Println("Path required")
      os.Exit(1)
    }

   wp_core_download(opts.Version, opts.Path)
  }
}
 
func main() {
  if len(os.Args) == 1 {
    fmt.Println("Command required")
    os.Exit(1)
  }

  command := os.Args[1]
  handle_command(command)
}