exec { "download-go":
	command => "wget http://golang.org/dl/go1.3.linux-amd64.tar.gz -O /opt/go.tar.gz",
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
    ensure   => "4.9.1.1-1ubuntu0.2",
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
	creates => [ "/dtop-dist/dist/dtop_0.3-linux-amd64.deb",
				 "/dtop-dist/dist/dtop_0.3-linux-amd64.rpm",
				 "/dtop-dist/dist/dtop_0.3-linux-i386.deb",
				 "/dtop-dist/dist/dtop_0.3-linux-i386.rpm" ],
	path => [ "/usr/bin/", "/usr/local/bin/", "/bin/" ]	
}
