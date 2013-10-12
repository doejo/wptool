package main
 
import(
  "fmt"
  "os"
  "os/exec"
  "bytes"
  "net/http"
  "io/ioutil"
  "regexp"
  "strings"
  "github.com/jessevdk/go-flags"
)

const(
  WP_DOWNLOAD_URL  = "http://wordpress.org/wordpress-%s.tar.gz"
  WP_VERSIONS_FILE = "https://raw.github.com/doejo/wptool/master/core-versions.txt"
)

type CoreConfigOptions struct {
  DbName string `long:"dbname" description:"Set the database name"`
  DbHost string `long:"dbhost" description:"Set the database host. Default: 'localhost'"`
  DbUser string `long:"dbuser" description:"Set the database user"`
  DbPass string `long:"dbpass" description:"Set the database password"`
}

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

func wp_core_version(path string) {
  version_path := fmt.Sprintf("%s/wp-includes/version.php", path)

  if !fileExists(version_path) {
    fmt.Println("Not a wordpress core")
    os.Exit(1)
  }

  buff, err := ioutil.ReadFile(version_path)
  if err != nil {
    fmt.Println("Unable to read version file")
    os.Exit(1)
  }

  exp := regexp.MustCompile(`wp_version = '(.*)'`)
  match := exp.FindString(string(buff))

  if len(match) == 0 {
    fmt.Println("Unable to find version")
    os.Exit(1)
  }

  chunks  := strings.Split(strings.TrimSpace(match), " ")
  version := strings.Replace(chunks[len(chunks) - 1], "'", "", -1)

  fmt.Println("Installed version:", version) 
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

func wp_core_download(version string, path string, force bool) {
  url := fmt.Sprintf(WP_DOWNLOAD_URL, version)
  temp := "/tmp/wordpress.tar.gz"

  if fileExists(temp) {
    run(fmt.Sprintf("rm -f %s", temp))
  }

  if fileExists(path) && !force {
    fmt.Println("Target path already exists!")
    os.Exit(1)
  }

  /* Download specified wordpress core */
  cmd := fmt.Sprintf("wget -O %s %s", temp, url)
  out, err := run(cmd)

  if err != nil {
    fmt.Println("Unable to download Wordpress Core!")
    fmt.Println(out)
    os.Exit(1)
  }
 
  /* Extract downloaded core tarball */
  out, err = run(fmt.Sprintf("cd /tmp && tar -zxf %s", temp))
  if err != nil {
    fmt.Println("Failed to extract wordpress core:", err)
    fmt.Println(out)
    os.Exit(1)
  }

  /* Remove existing directory with `--force` option */
  if fileExists(path) && force {
    run(fmt.Sprintf("rm -rf %s", path))
  }

  /* Move extracted files to target directory */
  _, err = run(fmt.Sprintf("mv /tmp/wordpress %s", path))
  if (err != nil) {
    fmt.Println("Failed to move extracted core")
    os.Exit(1)
  }

  /* Cleanup */
  run(fmt.Sprintf("rm -f %s", temp))

  /* Print installed version */
  wp_core_version(path)
}

func wp_core_config(path string) {

}

func handle_command(command string) {
  if command == "core:version" {
    var opts struct {
      Path string `short:"p" long:"path" description:"Path to core"`
    }

    _, err := flags.ParseArgs(&opts, os.Args)
    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    }

    if len(opts.Path) == 0 {
      opts.Path, _ = os.Getwd()
    }

    wp_core_version(opts.Path)
    return
  }

  if command == "core:list" {
    wp_core_list()
    return
  }

  if command == "core:download" {
    var opts struct {
      Version string `short:"v" long:"version" description:"Core version"`
      Path string    `short:"p" long:"path" description:"Path to install"`
      Force bool     `short:"f" long:"force" description:"Force override"`
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

   wp_core_download(opts.Version, opts.Path, opts.Force)
   return
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