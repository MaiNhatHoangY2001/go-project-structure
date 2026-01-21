package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig   `mapstructure:"server"`
	GRPC        GRPCConfig     `mapstructure:"grpc"`
	AuthService GRPCServiceConfig `mapstructure:"auth_service"`
	TodoService GRPCServiceConfig `mapstructure:"todo_service"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Logger      LoggerConfig   `mapstructure:"logger"`
	Environment string         `mapstructure:"environment"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type GRPCConfig struct {
	AuthService GRPCServiceConfig `mapstructure:"auth_service"`
	TodoService GRPCServiceConfig `mapstructure:"todo_service"`
}

type GRPCServiceConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

func LoadConfig(path ...string) (*Config, error) {
	// If path is provided, use it; otherwise use default
	configPath := "configs/app.yaml"
	if len(path) > 0 && path[0] != "" {
		configPath = path[0]
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set default values
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("grpc.auth_service.port", "50051")
	viper.SetDefault("grpc.auth_service.host", "localhost")
	viper.SetDefault("grpc.todo_service.port", "50052")
	viper.SetDefault("grpc.todo_service.host", "localhost")
	viper.SetDefault("auth_service.port", "50051")
	viper.SetDefault("auth_service.host", "localhost")
	viper.SetDefault("todo_service.port", "50052")
	viper.SetDefault("todo_service.host", "localhost")
	viper.SetDefault("database.mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.mongodb.database", "todoapp")
	viper.SetDefault("database.mongodb.timeout", 10)
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration", 24)
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.encoding", "json")

	// Read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v. Using defaults.", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
