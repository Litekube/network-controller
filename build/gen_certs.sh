#!/bin/bash

# demo: ./gen_certs.sh $ip

ip=$1

mkdir -p ../certs/init/test1
mkdir -p ../certs/init/test2

# generate network certs
echo 'start generate network certs'
cd ../certs/init/test1

cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "expiry": "8760h",
        "usages": [
          "signing",
          "key encipherment",
          "server auth"
        ]
      },
      "client": {
        "expiry": "8760h",
        "usages": [
          "signing",
          "key encipherment",
          "client auth"
        ]
      }
    }
  }
}
EOF

cat > ca-csr.json <<EOF
{
  "CN": "network-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "Beijing",
      "ST": "Beijing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

cat > server-csr.json <<EOF
{
  "CN": "network-server",
  "hosts": [
    "127.0.0.1",
    "10.1.1.1",
    "$ip"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "BeiJing",
      "ST": "BeiJing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

cat > client-csr.json <<EOF
{
  "CN": "network-client",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "BeiJing",
      "ST": "BeiJing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

# 生成ca证书和私钥 ca-key.pem ca.pem
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

# 生成network server的证书和私钥 server-key.pem server.pem
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server

# # 生成network client的证书和私钥 server-key.pem server.pem
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client-csr.json | cfssljson -bare client


# generate grpc certs
echo 'start generate network certs'
cd ../test2

cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "server": {
        "expiry": "8760h",
        "usages": [
          "signing",
          "key encipherment",
          "server auth"
        ]
      },
      "client": {
        "expiry": "8760h",
        "usages": [
          "signing",
          "key encipherment",
          "client auth"
        ]
      }
    }
  }
}
EOF

cat > ca-csr.json <<EOF
{
  "CN": "grpc-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "Beijing",
      "ST": "Beijing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

cat > server-csr.json <<EOF
{
  "CN": "grpc-server",
  "hosts": [
    "127.0.0.1",
    "10.1.1.1",
    "$ip"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "BeiJing",
      "ST": "BeiJing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

cat > client-csr.json <<EOF
{
  "CN": "grpc-client",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "L": "BeiJing",
      "ST": "BeiJing",
      "O": "litekube",
      "OU": "network-controller"
    }
  ]
}
EOF

# 生成ca证书和私钥 ca-key.pem ca.pem
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

# 生成grpc server的证书和私钥 server-key.pem server.pem
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server

# # 生成grpc client的证书和私钥 server-key.pem server.pem
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client-csr.json | cfssljson -bare client