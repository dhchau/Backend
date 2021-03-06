# Install postgresql 9.4
class postgres {
  $ppa_source = 'anonscm.debian.org/loggerhead/pkg-postgresql/postgresql-common/trunk/download/head:/apt.postgresql.org.s-20130224224205-px3qyst90b3xp8zj-1/apt.postgresql.org.sh'

  # Install postgresql ubuntu ppa
  exec {'installppa':
    command => "curl -ssL ${ppa_source} | bash -x",
    unless  => 'ls /etc/apt/sources.list.d/pgdg.list',
  }

  # Enable postgres updates to 9.4
  exec {
    'use-9.4':
      command => "sed -i '/^[^#]/ s/$/ 9.4/' /etc/apt/sources.list.d/pgdg.list",
      require => Exec['installppa'];
    'update-apt':
      command => 'apt-get update',
      require => Exec['use-9.4']
  }

  # Install postgres 9.4
  package {'postgresql-9.4':
    ensure  => installed,
    require => Exec['update-apt'],
  }

  # Set up postgres configuration
  file {'postgresql.conf':
    ensure  => file,
    path    => '/etc/postgresql/9.4/main/postgresql.conf',
    source  => 'puppet:///modules/postgres/postgresql.conf',
    require => Package['postgresql-9.4'],
  }

  # Ensure service starts at boot
  service {'postgresql':
    ensure    => running,
    enable    => true,
    subscribe => File['postgresql.conf'],
  }
}
