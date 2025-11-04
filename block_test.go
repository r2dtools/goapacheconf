package goapacheconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockIfModules(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	require.Equal(t, []string{"ssl"}, vBlock.IfModules)
	blocks := vBlock.FindBlocks(Proxy)
	require.Len(t, blocks, 1)

	block := blocks[0]
	require.Equal(t, []string{"ssl", "proxy_http"}, block.IfModules)
}

func TestFindIfModuleBlocks(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[1]
	blocks := vBlock.FindIfModuleBlocks()
	require.Len(t, blocks, 2)

	iBlock := blocks[0]
	blocks = iBlock.FindIfModuleBlocksByModuleName("rewrite")
	require.Len(t, blocks, 1)
	require.Equal(t, []string{"proxy_http"}, blocks[0].IfModules)

	blocks = vBlock.FindIfModuleBlocksByModuleName("rewrite")
	require.Len(t, blocks, 1)
}

func TestAddDirectiveToBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		configFile := getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks := configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)

		vBlock := vBlocks[0]
		directive := NewDirective("Test", []string{"test"})
		directive.AppendNewLine()
		directive = vBlock.PrependDirective(directive)
		_, err := configFile.Dump()
		require.Nil(t, err)

		configFile = getConfigFile(t, r2dtoolsConfigFileName)
		vBlocks = configFile.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
		require.Len(t, vBlocks, 2)
		vBlock = vBlocks[0]

		directives := vBlock.FindDirectives("Test")
		require.Len(t, directives, 1)
		require.Equal(t, []string{"test"}, directives[0].GetValues())
	})
}

func TestDeleteDirectiveFromBlock(t *testing.T) {
	testWithConfigFileRollback(t, r2dtoolsConfigFilePath, func(t *testing.T) {
		config, block := getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		block.DeleteDirectiveByName(UseCanonicalName)
		err := config.Dump()
		require.Nil(t, err)

		_, block = getFirstVirtualHostBlock(t, "r2dtools.work.gd")
		directives := block.FindDirectives(UseCanonicalName)
		require.Empty(t, directives)
	})
}

func TestGetOrder(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[0]
	directives := vBlock.FindAlliasDirectives()
	require.NotEmpty(t, directives)

	directive := directives[0]
	order := vBlock.GetDirectiveOrder(directive.Directive)
	require.Equal(t, 13, order)

	blocks := vBlock.FindDirectoryBlocks()
	require.NotEmpty(t, blocks)
	block := blocks[0]

	order = vBlock.GetBlockOrder(block.Block)
	require.Equal(t, 24, order)
}

func TestChangeDirectoryOrder(t *testing.T) {
	config := getConfig(t)

	vBlocks := config.FindVirtualHostBlocksByServerName("r2dtools.work.gd")
	require.Len(t, vBlocks, 2)

	vBlock := vBlocks[1]

	directives := vBlock.FindDirectives(CustomLog)
	require.Len(t, directives, 1)
	directive := directives[0]

	order := vBlock.GetDirectiveOrder(directive)
	require.Equal(t, 4, order)

	ifModuleBlocks := vBlock.FindIfModuleBlocks()
	require.Len(t, ifModuleBlocks, 2)
	ifModuleBlock := ifModuleBlocks[1]

	nOrder := vBlock.GetBlockOrder(ifModuleBlock.Block)
	require.Equal(t, 9, nOrder)

	vBlock.ChangeDirectiveOrder(directive, nOrder)
	content := vBlock.Dump()

	exepcted := `<VirtualHost 127.0.0.1:7080 >
	ServerName "r2dtools.work.gd"
	ServerAlias "www.r2dtools.work.gd"
	ServerAlias "ipv4.r2dtools.work.gd"
	UseCanonicalName Off

	ErrorLog "/var/www/vhosts/system/r2dtools.work.gd/logs/error_log"

	# mailconfig
	<IfModule mod_proxy_http.c >
		<IfModule mod_rewrite.c >
			RewriteEngine On
			RewriteCond %{REQUEST_URI} ^/autodiscover/autodiscover\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^(/\.well-known/autoconfig)?/mail/config\-v1\.1\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^/email\.mobileconfig$ [NC]
			RewriteRule ^(.*)$ http://127.0.0.1:8880/mailconfig/ [P,QSA,L,E=REQUEST_URI:%{REQUEST_URI},E=HOST:%{HTTP_HOST}]
		</IfModule>
		<Proxy "http://127.0.0.1:8880/mailconfig/" >
			RequestHeader set X-Host "%{HOST}e"
			RequestHeader set X-Request-URI "%{REQUEST_URI}e"
		</Proxy>
	</IfModule>
	# mailconfig

	CustomLog /var/www/vhosts/system/r2dtools.work.gd/logs/access_log plesklog
	<IfModule mod_rewrite.c >
		RewriteEngine On
		RewriteCond %{HTTPS} off
		RewriteRule ^ https://%{HTTP_HOST}%{REQUEST_URI} [R=301,L,QSA]
	</IfModule>

	Alias /.well-known/acme-challenge "/var/www/vhosts/default/htdocs/.well-known/acme-challenge"
</VirtualHost>`

	require.Equal(t, exepcted, content)

	vBlock.ChangeDirectiveOrder(directive, 1)
	content = vBlock.Dump()

	exepcted = `<VirtualHost 127.0.0.1:7080 >
	ServerName "r2dtools.work.gd"
	CustomLog /var/www/vhosts/system/r2dtools.work.gd/logs/access_log plesklog
	ServerAlias "www.r2dtools.work.gd"
	ServerAlias "ipv4.r2dtools.work.gd"
	UseCanonicalName Off

	ErrorLog "/var/www/vhosts/system/r2dtools.work.gd/logs/error_log"

	# mailconfig
	<IfModule mod_proxy_http.c >
		<IfModule mod_rewrite.c >
			RewriteEngine On
			RewriteCond %{REQUEST_URI} ^/autodiscover/autodiscover\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^(/\.well-known/autoconfig)?/mail/config\-v1\.1\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^/email\.mobileconfig$ [NC]
			RewriteRule ^(.*)$ http://127.0.0.1:8880/mailconfig/ [P,QSA,L,E=REQUEST_URI:%{REQUEST_URI},E=HOST:%{HTTP_HOST}]
		</IfModule>
		<Proxy "http://127.0.0.1:8880/mailconfig/" >
			RequestHeader set X-Host "%{HOST}e"
			RequestHeader set X-Request-URI "%{REQUEST_URI}e"
		</Proxy>
	</IfModule>
	# mailconfig

	<IfModule mod_rewrite.c >
		RewriteEngine On
		RewriteCond %{HTTPS} off
		RewriteRule ^ https://%{HTTP_HOST}%{REQUEST_URI} [R=301,L,QSA]
	</IfModule>

	Alias /.well-known/acme-challenge "/var/www/vhosts/default/htdocs/.well-known/acme-challenge"
</VirtualHost>`

	require.Equal(t, exepcted, content)

	aliasDirectives := vBlock.FindAlliasDirectives()
	require.Len(t, aliasDirectives, 1)
	aliasDirective := aliasDirectives[0]

	vBlock.ChangeDirectiveOrder(aliasDirective.Directive, 1)
	content = vBlock.Dump()

	exepcted = `<VirtualHost 127.0.0.1:7080 >
	ServerName "r2dtools.work.gd"
	Alias /.well-known/acme-challenge "/var/www/vhosts/default/htdocs/.well-known/acme-challenge"
	CustomLog /var/www/vhosts/system/r2dtools.work.gd/logs/access_log plesklog
	ServerAlias "www.r2dtools.work.gd"
	ServerAlias "ipv4.r2dtools.work.gd"
	UseCanonicalName Off

	ErrorLog "/var/www/vhosts/system/r2dtools.work.gd/logs/error_log"

	# mailconfig
	<IfModule mod_proxy_http.c >
		<IfModule mod_rewrite.c >
			RewriteEngine On
			RewriteCond %{REQUEST_URI} ^/autodiscover/autodiscover\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^(/\.well-known/autoconfig)?/mail/config\-v1\.1\.xml$ [NC,OR]
			RewriteCond %{REQUEST_URI} ^/email\.mobileconfig$ [NC]
			RewriteRule ^(.*)$ http://127.0.0.1:8880/mailconfig/ [P,QSA,L,E=REQUEST_URI:%{REQUEST_URI},E=HOST:%{HTTP_HOST}]
		</IfModule>
		<Proxy "http://127.0.0.1:8880/mailconfig/" >
			RequestHeader set X-Host "%{HOST}e"
			RequestHeader set X-Request-URI "%{REQUEST_URI}e"
		</Proxy>
	</IfModule>
	# mailconfig

	<IfModule mod_rewrite.c >
		RewriteEngine On
		RewriteCond %{HTTPS} off
		RewriteRule ^ https://%{HTTP_HOST}%{REQUEST_URI} [R=301,L,QSA]
	</IfModule>

</VirtualHost>`

	require.Equal(t, exepcted, content)
}
