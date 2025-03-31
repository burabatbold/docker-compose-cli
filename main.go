package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Image         string            `yaml:"image"`
	Container_name string           `yaml:"container_name"`
	Environment   []string          `yaml:"environment,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
}

type DockerCompose struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type ServiceConfig struct {
	name     string
	username string
	password string
	port     string
}

func main() {
	services := []string{"minio", "postgresql", "mssql", "mysql"}
	selectedServices := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select services to include:",
		Options: services,
	}
	survey.AskOne(prompt, &selectedServices)

	dockerCompose := DockerCompose{
		Version:  "3.8",
		Services: make(map[string]Service),
	}

	for _, serviceName := range selectedServices {
		config := getServiceConfig(serviceName)
		service := createService(config)
		dockerCompose.Services[serviceName] = service
	}

	yamlData, err := yaml.Marshal(&dockerCompose)
	if err != nil {
		fmt.Printf("Error marshaling YAML: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile("docker-compose.yml", yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Docker Compose file generated successfully!")
}

func getServiceConfig(service string) ServiceConfig {
	config := ServiceConfig{name: service}
	
	questions := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Enter username for %s:", service),
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: fmt.Sprintf("Enter password for %s:", service),
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Enter port for %s:", service),
			},
		},
	}

	answers := struct {
		Username string
		Password string
		Port     string
	}{}

	survey.Ask(questions, &answers)

	config.username = answers.Username
	config.password = answers.Password
	config.port = answers.Port

	return config
}

func createService(config ServiceConfig) Service {
	service := Service{
		Container_name: fmt.Sprintf("%s-container", config.name),
		Ports:         []string{fmt.Sprintf("%s:%s", config.port, getDefaultPort(config.name))},
	}

	switch config.name {
	case "minio":
		service.Image = "minio/minio:latest"
		service.Environment = []string{
			fmt.Sprintf("MINIO_ROOT_USER=%s", config.username),
			fmt.Sprintf("MINIO_ROOT_PASSWORD=%s", config.password),
		}
		service.Volumes = []string{"./minio/data:/data"}
		service.Ports = append(service.Ports, "9001:9001")
		
	case "postgresql":
		service.Image = "postgres:latest"
		service.Environment = []string{
			fmt.Sprintf("POSTGRES_USER=%s", config.username),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", config.password),
			"POSTGRES_DB=mydatabase",
		}
		service.Volumes = []string{"./postgresql/data:/var/lib/postgresql/data"}

	case "mssql":
		service.Image = "mcr.microsoft.com/mssql/server:2022-latest"
		service.Environment = []string{
			"ACCEPT_EULA=Y",
			fmt.Sprintf("SA_PASSWORD=%s", config.password),
			"MSSQL_PID=Express",
		}
		service.Volumes = []string{"./mssql/data:/var/opt/mssql/data"}

	case "mysql":
		service.Image = "mysql:latest"
		service.Environment = []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", config.password),
			fmt.Sprintf("MYSQL_USER=%s", config.username),
			fmt.Sprintf("MYSQL_PASSWORD=%s", config.password),
			"MYSQL_DATABASE=mydatabase",
		}
		service.Volumes = []string{"./mysql/data:/var/lib/mysql"}
	}

	return service
}

func getDefaultPort(service string) string {
	switch service {
	case "minio":
		return "9000"
	case "postgresql":
		return "5432"
	case "mssql":
		return "1433"
	case "mysql":
		return "3306"
	default:
		return "80"
	}
} 