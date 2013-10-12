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
  "github.com/hoisie/mustache"
)

const VERSION = "0.1.0"

const(
  WP_DOWNLOAD_URL  = "http://wordpress.org/wordpress-%s.tar.gz"
  WP_VERSIONS_FILE = "https://raw.github.com/doejo/wptool/master/core-versions.txt"
  WP_CONFIG_FILE   = "https://raw.github.com/doejo/wptool/master/templates/wp-config.mustache"
  WP_SALTS_API     = "https://api.wordpress.org/secret-key/1.1/salt/"
)

type CoreConfigOptions struct {
  Path      string `short:"p" long:"path" description:"Path to wordpress core"`
  Template  string `short:"t" long:"template" description:"Config template"`
  Force     bool   `short:"f" long:"force" description:"Force config update"`
  DbName    string `long:"dbname" description:"Set the database name"`
  DbHost    string `long:"dbhost" description:"Set the database host. Default: 'localhost'"`
  DbUser    string `long:"dbuser" description:"Set the database user"`
  DbPass    string `long:"dbpass" description:"Set the database password"`
  DbCharset string `long:"dbcharset" description:"Set the database charset"`
  DbCollate string `long:"dbcollate" description: "Set the database collate type"`
  DbPrefix  string `long:"dbprefix" description:"Set the database prefix"`
}

type CoreDownloadOptions struct {
  Version string `short:"v" long:"version" description:"Core version"`
  Path    string `short:"p" long:"path" description:"Path to install"`
  Force   bool   `short:"f" long:"force" description:"Force override"`
}

func getUrlContents(url string) string {
  var err error
  var resp *http.Response
  var body []byte

  resp, err = http.Get(url)
  if err != nil {
    return ""
  }

  defer resp.Body.Close()

  body, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    return ""
  }

  return string(body)
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

func checkDbConnection(config *CoreConfigOptions) {
  _, err := run("which mysql")
  if err != nil {
    fmt.Println("MySQL client is not installed")
    os.Exit(1)
  }

  cmd := fmt.Sprintf(
    "mysql --no-defaults -h %s -u %s -p%s -e ';'",
    config.DbHost, config.DbUser, config.DbPass,
  )

  _, err = run(cmd)
  if err != nil {
    fmt.Println("Unable to establish connection to MySQL")
    os.Exit(1)
  }
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
  result := getUrlContents(WP_VERSIONS_FILE)

  if len(result) == 0 {
    fmt.Println("Unable to get list of versions")
    os.Exit(1)
  }

  fmt.Println(result)
}

func wp_core_download(options *CoreDownloadOptions) {
  url := fmt.Sprintf(WP_DOWNLOAD_URL, options.Version)
  temp := "/tmp/wordpress.tar.gz"

  if fileExists(options.Path) && !options.Force {
    fmt.Println("Path already exists")
    os.Exit(1)
  }

  /* Remove downloaded file if exists */
  if fileExists(temp) {
    fmt.Println("Removing temporary archive")
    run(fmt.Sprintf("rm -f %s", temp))
  }

  /* Download specified wordpress core */
  fmt.Println("Downloading core:", url)
  out, err := run(fmt.Sprintf("wget -O %s %s", temp, url))

  if err != nil {
    fmt.Println("Unable to download wordpress core")
    os.Exit(1)
  }
 
  /* Extract downloaded core tarball */
  out, err = run(fmt.Sprintf("cd /tmp && tar -zxf %s", temp))
  if err != nil {
    fmt.Println("Unable to extract core", err)
    fmt.Println(out)
    os.Exit(1)
  }

  /* Remove existing directory with `--force` option */
  if fileExists(options.Path) && options.Force {
    fmt.Println("Removing existing core")
    run(fmt.Sprintf("rm -rf %s", options.Path))
  }

  /* Move extracted files to target directory */
  fmt.Println("Extracting core")
  _, err = run(fmt.Sprintf("mv /tmp/wordpress %s", options.Path))
  if (err != nil) {
    fmt.Println("Failed to move extracted core")
    os.Exit(1)
  }

  /* Cleanup */
  run(fmt.Sprintf("rm -f %s", temp))

  /* Print installed version */
  wp_core_version(options.Path)
}

func wp_core_config(options *CoreConfigOptions) {
  var err error
  
  config_path := fmt.Sprintf("%s/wp-config.php", options.Path)

  if fileExists(config_path) && !options.Force {
    fmt.Println("Config file already exists")
    os.Exit(1)
  }

  /* Get keys and salts from wp api */
  salts := getUrlContents(WP_SALTS_API)
  if len(salts) == 0 {
    fmt.Println("Unable to get salts from wordpress API")
    os.Exit(1)
  }

  /* Check database connection */
  checkDbConnection(options)

  /* Setup config for rendering */
  config := map[string]string {
    "dbname":    options.DbName,
    "dbuser":    options.DbUser,
    "dbpass":    options.DbPass,
    "dbhost":    options.DbHost,
    "dbcharset": options.DbCharset,
    "dbcollate": options.DbCollate,
    "dbprefix":  options.DbPrefix,
    "keys-and-salts": salts,
  }

  /* Remove existing config file */
  run(fmt.Sprintf("rm -f %s", config_path))

  result := mustache.RenderFile(options.Template, config)
  err = ioutil.WriteFile(config_path, []byte(result), 0644)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  /* Print result */
  fmt.Printf("New config file generated at %s\n", config_path)
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
    opts := CoreDownloadOptions {}

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

    wp_core_download(&opts)
    return
  }

  if command == "core:config" {
    opts := CoreConfigOptions {}
    
    _, err := flags.ParseArgs(&opts, os.Args)
    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    }

    if len(opts.Path) == 0 {
      opts.Path, _ = os.Getwd()
    }

    if len(opts.DbHost) == 0 {
      opts.DbHost = "localhost"
    }

    if len(opts.DbPrefix) == 0 {
      opts.DbPrefix = "wp_"
    }

    if len(opts.DbCharset) == 0 {
      opts.DbCharset = "utf8"
    }

    if len(opts.DbName) == 0 {
      fmt.Println("Database name required")
      os.Exit(1)
    }

    if len(opts.DbUser) == 0 {
      fmt.Println("Database user required")
      os.Exit(1)
    }

    if len(opts.DbPass) == 0 {
      fmt.Println("Database password required")
      os.Exit(1)
    }

    if len(opts.Template) == 0 {
      opts.Template = getUrlContents(WP_CONFIG_FILE);

      if len(opts.Template) == 0 {
        fmt.Println("Config template path required")
        os.Exit(1)
      }
    }

    wp_core_config(&opts)
    return
  }

  if command == "version" {
    fmt.Printf("wptool v%s\n", VERSION)
    return
  }

  fmt.Println("Invalid command")
  os.Exit(1)
}
 
func main() {
  if len(os.Args) == 1 {
    fmt.Println("Command required")
    os.Exit(1)
  }

  command := os.Args[1]
  handle_command(command)
}