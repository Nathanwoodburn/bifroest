package configuration

import (
	"gopkg.in/yaml.v3"

	"github.com/engity-com/bifroest/pkg/template"
)

var (
	DefaultEnvironmentRemoteLoginAllowed = template.BoolOf(true)

	DefaultEnvironmentRemoteHost     = template.MustNewString("")
	DefaultEnvironmentRemotePort     = template.MustNewString("22")
	DefaultEnvironmentRemoteUser     = template.MustNewString("{{.authorization.user.name}}")
	DefaultEnvironmentRemotePassword = template.MustNewString("")

	DefaultEnvironmentRemoteShellCommand = template.MustNewStrings()
	DefaultEnvironmentRemoteExecCommand  = template.MustNewStrings()
	DefaultEnvironmentRemoteSftpCommand  = template.MustNewStrings()

	DefaultEnvironmentRemoteBanner                = template.MustNewString("")
	DefaultEnvironmentRemotePortForwardingAllowed = template.BoolOf(true)

	_ = RegisterEnvironmentV(func() EnvironmentV {
		return &EnvironmentRemote{}
	})
)

type EnvironmentRemote struct {
	LoginAllowed template.Bool `yaml:"loginAllowed,omitempty"`

	Host     template.String  `yaml:"host"`
	Port     template.String  `yaml:"port,omitempty"`
	User     template.String  `yaml:"user,omitempty"`
	Password template.String  `yaml:"password,omitempty"`

	ShellCommand template.Strings `yaml:"shellCommand,omitempty"`
	ExecCommand  template.Strings `yaml:"execCommand,omitempty"`
	SftpCommand  template.Strings `yaml:"sftpCommand,omitempty"`

	Banner                template.String `yaml:"banner,omitempty"`
	PortForwardingAllowed template.Bool   `yaml:"portForwardingAllowed,omitempty"`
}

func (this *EnvironmentRemote) SetDefaults() error {
	return setDefaults(this,
		fixedDefault("loginAllowed", func(v *EnvironmentRemote) *template.Bool { return &v.LoginAllowed }, DefaultEnvironmentRemoteLoginAllowed),

		fixedDefault("host", func(v *EnvironmentRemote) *template.String { return &v.Host }, DefaultEnvironmentRemoteHost),
		fixedDefault("port", func(v *EnvironmentRemote) *template.String { return &v.Port }, DefaultEnvironmentRemotePort),
		fixedDefault("user", func(v *EnvironmentRemote) *template.String { return &v.User }, DefaultEnvironmentRemoteUser),
		fixedDefault("password", func(v *EnvironmentRemote) *template.String { return &v.Password }, DefaultEnvironmentRemotePassword),

		fixedDefault("shellCommand", func(v *EnvironmentRemote) *template.Strings { return &v.ShellCommand }, DefaultEnvironmentRemoteShellCommand),
		fixedDefault("execCommand", func(v *EnvironmentRemote) *template.Strings { return &v.ExecCommand }, DefaultEnvironmentRemoteExecCommand),
		fixedDefault("sftpCommand", func(v *EnvironmentRemote) *template.Strings { return &v.SftpCommand }, DefaultEnvironmentRemoteSftpCommand),

		fixedDefault("banner", func(v *EnvironmentRemote) *template.String { return &v.Banner }, DefaultEnvironmentRemoteBanner),
		fixedDefault("portForwardingAllowed", func(v *EnvironmentRemote) *template.Bool { return &v.PortForwardingAllowed }, DefaultEnvironmentRemotePortForwardingAllowed),
	)
}

func (this *EnvironmentRemote) Trim() error {
	return trim(this,
		noopTrim[EnvironmentRemote]("loginAllowed"),

		noopTrim[EnvironmentRemote]("host"),
		noopTrim[EnvironmentRemote]("port"),
		noopTrim[EnvironmentRemote]("user"),
		noopTrim[EnvironmentRemote]("password"),

		noopTrim[EnvironmentRemote]("shellCommand"),
		noopTrim[EnvironmentRemote]("execCommand"),
		noopTrim[EnvironmentRemote]("sftpCommand"),

		noopTrim[EnvironmentRemote]("banner"),
		noopTrim[EnvironmentRemote]("portForwardingAllowed"),
	)
}

func (this *EnvironmentRemote) Validate() error {
	return validate(this,
		func(v *EnvironmentRemote) (string, validator) { return "loginAllowed", &v.LoginAllowed },

		func(v *EnvironmentRemote) (string, validator) { return "host", &v.Host },
		notZeroValidate("host", func(v *EnvironmentRemote) *template.String { return &v.Host }),
		func(v *EnvironmentRemote) (string, validator) { return "port", &v.Port },
		func(v *EnvironmentRemote) (string, validator) { return "user", &v.User },
		func(v *EnvironmentRemote) (string, validator) { return "password", &v.Password },

		func(v *EnvironmentRemote) (string, validator) { return "shellCommand", &v.ShellCommand },
		func(v *EnvironmentRemote) (string, validator) { return "execCommand", &v.ExecCommand },
		func(v *EnvironmentRemote) (string, validator) { return "sftpCommand", &v.SftpCommand },

		func(v *EnvironmentRemote) (string, validator) { return "banner", &v.Banner },
		func(v *EnvironmentRemote) (string, validator) { return "portForwardingAllowed", &v.PortForwardingAllowed },
	)
}

func (this *EnvironmentRemote) UnmarshalYAML(node *yaml.Node) error {
	return unmarshalYAML(this, node, func(target *EnvironmentRemote, node *yaml.Node) error {
		type raw EnvironmentRemote
		return node.Decode((*raw)(target))
	})
}

func (this EnvironmentRemote) IsEqualTo(other any) bool {
	if other == nil {
		return false
	}
	switch v := other.(type) {
	case EnvironmentRemote:
		return this.isEqualTo(&v)
	case *EnvironmentRemote:
		return this.isEqualTo(v)
	default:
		return false
	}
}

func (this EnvironmentRemote) isEqualTo(other *EnvironmentRemote) bool {
	return isEqual(&this.LoginAllowed, &other.LoginAllowed) &&
		isEqual(&this.Host, &other.Host) &&
		isEqual(&this.Port, &other.Port) &&
		isEqual(&this.User, &other.User) &&
		isEqual(&this.Password, &other.Password) &&
		isEqual(&this.ShellCommand, &other.ShellCommand) &&
		isEqual(&this.ExecCommand, &other.ExecCommand) &&
		isEqual(&this.SftpCommand, &other.SftpCommand) &&
		isEqual(&this.Banner, &other.Banner) &&
		isEqual(&this.PortForwardingAllowed, &other.PortForwardingAllowed)
}

func (this EnvironmentRemote) Types() []string {
	return []string{"remote"}
}

func (this EnvironmentRemote) FeatureFlags() []string {
	return []string{"remote"}
}