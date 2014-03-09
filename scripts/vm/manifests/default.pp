exec { "download-go":
	command => "wget https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz -O go.tar.gz",
	creates => "/opt/go.tar.gz",
	cwd => "/opt/",
	path => [ "/usr/bin/", "/bin/" ]
}
->
exec { "unpack-go":
	command => "tar xzf go.tar.gz",
	creates => "/opt/go/bin/go",
	cwd => "/opt/",
	path => [ "/usr/bin/", "/bin/" ]
}
->
package { "rpm":
    ensure   => "installed",
}
->
package { "fpm":
    ensure   => "1.0.2",
    provider => "gem",
}
->
exec {  "make-distros":
	command => "make dist-all",
	cwd => "/dtop-dist",
	creates => [ "/dtop-dist/dist/dtop_0.1-linux-amd64.deb",
				 "/dtop-dist/dist/dtop_0.1-linux-amd64.rpm" ],
	path => [ "/usr/bin/", "/usr/local/bin/", "/bin/" ]	
}