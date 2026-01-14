package main

// need you add the cert creation lib from mem-sdk
import (
	"log"
	"os"

	"github.com/odio4u/mem-sdk/certengine/pkg"
	"gopkg.in/yaml.v3"
)

type Seeder struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
	Dns  string `yaml:"dns"`
	Name string `yaml:"name"`
}

type Config struct {
	Version string `yaml:"version"`
	Seeder  Seeder `yaml:"Seeder"`
}

var YamlConfig *Config

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	routerIps := []string{YamlConfig.Seeder.IP}
	dns := []string{YamlConfig.Seeder.Dns}

	_, err := pkg.GenerateSelfSignedGPR(YamlConfig.Seeder.Name, routerIps, dns)
	pkg.Must(err)

	log.Println("All operations completed successfully.")
}

func init() {
	data, err := os.ReadFile("seeder-config.yaml")
	log.Println("load the config")

	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error parsing YAML file: %v", err)
	}

	YamlConfig = &config

}
