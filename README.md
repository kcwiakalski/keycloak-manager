# Keycloak Manager - handle you Keycloak server configuration with ease

Table of Content
- [What is it](#what-is-it?)
- [Motivation](#motivation)
- [How it works?](#how-it-works?)
- [Usage](#usage)
   - [Basic usage](#basic-usage)
   - [Available commands](#available-commands)
- [Build](#build)
- [Features](#features)
   - [Supported operations for Keycloak client's](supported-operations-for-keycloak-client's)
   - [Configuration data model](#configuration-data-model)
- [Future development](#future-development)
- [Configuration examples and generated diff json](#configuration-examples-and-generated-diff-in-json)
   - [Not existing client](#not-existing-client)
   - [Adding scope and new resource to existing client](#adding-scope-and-new-resource-to-existing-client)
   - [Removing policy and permission](#removing-policy-and-permission)

# What is it?
Tool to configure and manage Keycloak server in terms of clients defined for realm. It is able to take configuration for client, check it against live Keycloak instance and determine what needs to be done in order to align its configuration to desired state.

# Motivation
Working with Keyclock is mostly focus on configuring its various aspects. Usually we work with some local or test instance, manully configure required things and then want to apply it on some other instance. Very often is it the the same insance but with different realm (e.g. when we share Keycloak server between various version of the same service). Keycloak provides export/import capabilities, but unfornutelly exporting realm and clients settings comes with already configured IDs for configurations elements, which renders applying it twice on same instance impossible. Purpose of this tool is to be able to apply desired state on living server and align its configuration with developer's expectation.

# How it works?
Appling configuration is two-step process:
1. Curent state evaluation
- take desired state writen as data in json format
- connect to Keycloak instance 
- cross-check its state with expectations expresed via json structure
- generate diff (in json format) with operations that need to be done to align reality with expecetation
2. Change application
- take generated diff in json format
- connect to Keycloak instance
- execute operation defined in diff

# Usage
To see all possible parameters combination use --help flag on main program.
## Basic usage
```bash
$ ./keycloak-manager --help
Usage: keycloak-manager <command>

Flags:
  -h, --help            Show context-sensitive help.
  -p, --port=INT        Port on which Keycloak Admin Api is available
  -h, --host=STRING     Server hosting Keycloak instalation
  -u, --user=STRING     Username with administrative rights
      --pass=STRING     Password for user with administrative rights. It is highly discuraged to use this flag directly
  -r, --realm=STRING    Realm holding user with administrative rights, usually the same as realm that is target for operation

Commands:
  client
    Operates on client configuration

Run "keycloak-manager <command> --help" for more information on a command.
```
## Available commands
* Client handling:
```bash
$ ./keycloak-manager client --help
Usage: keycloak-manager client

Operates on client configuration

Flags:
  -h, --help                                  Show context-sensitive help.
  -p, --port=INT                              Port on which Keycloak Admin Api is available
  -h, --host=STRING                           Server hosting Keycloak instalation
  -u, --user=STRING                           Username with administrative rights
      --pass=STRING                           Password for user with administrative rights. It is highly discuraged to use this flag directly
  -r, --realm=STRING                          Realm holding user with administrative rights, usually the same as realm that is target for operation

  -f, --file="client-config.json"             Path to file with client configuration
  -m, --mode="diff"                           Indicates what should be done with config file
  -o, --output="client-config-change.json"    For diff flag indicates file name what will hold operations to apply
```

# Features
## Supported operations for Keycloak client's
| Context    | Add | Delete | Update |
|------------|-----|--------|--------|
|client      |  x  |        |        |
|scopes      |  x  |    x   |        | 
|resources   |  x  |    x   |        | 
|policies    |  x  |    x   |        | 
|permissions |  x  |    x   |        | 
|realm groups|  x  |        |        |
## Configuration data model
Configration is describe by json object containing attributes for specific aspect of client configuration. Configuration uses standard Keycloak data structures defined for its [REST Api](https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_definitions). Basic json looks like this:
```json
{
    "$schema": "http://json-schema.org/schema#",
    "properties": {
        "clientDefinition": {
            "type": "object",
            "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_clientrepresentation"
        },
        "groups": {
            "type": "array",
            "items": {
                "type": "object",
                "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_grouprepresentation"
            }
        },
        "scopes": {
            "type": "array",
            "items": {
                "type": "object",
                "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_clientscoperepresentation"
            }
        },
        "resources": {
            "type": "array",
            "items": {
                "type": "object",
                "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_resourcerepresentation"
            }
        },
        "policies": {
            "type": "array",
            "items": {
                "type": "object",
                "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_policyrepresentation"
            }
        },
        "permissions": {
            "type": "array",
            "items": {
                "type": "object",
                "description": "https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_policyrepresentation"
            }
        }
    }
}
```

# Build
There are two way you can build this tool on your own:
1. If you have golang local setup ready just clone repository and call
   ```
   go build
   ```
2. Alternatively if you have Docker installed on your machine you can use official Golang image and build biniaries without setting up local development environment. 
   - build for linux 64-bit 
      ```
      docker run --rm -v "$PWD":/usr/src/keycloak-manager -w /usr/src/keycloak-manager -e GOOS=linux -e GOARCH=amd64 golang:1.14 go build
      ```
   - build for windows 64-bit
      ```
      docker run --rm -v "$PWD":/usr/src/keycloak-manager -w /usr/src/keycloak-manager -e GOOS=windows -e GOARCH=amd64 golang:1.14 go build
      ```

# Future development
- client
   - handle update operations
   - manage client roles
   ...
- realm - new configuration area
   - basic realm management
   - identiy provider declaration
   - email configuration
   ...

# Configuration examples and generated diff in json
## Not existing client
Let's assume we empty realm named *sample-realm* and we want to create client named *sample-client*, whose configuration is as follows ([also available here](examples/data-access-client-conf.json)):
```json
{
    "clientDefinition": {
        "enabled": true,
        "redirectUris": [
            "*"
        ],
        "clientId": "sample-client",
        "protocol": "openid-connect",
        "attributes": {
            "login_theme": "base"
        },
        "authorizationServicesEnabled": true,
        "serviceAccountsEnabled": true,
        "publicClient": false
    },
    "groups": [
        {
            "name": "data-user",
            "path": "/data-user"
        },
        {
            "name": "data-creator",
            "path": "/data-user/data-creator"
        },
        {
            "name": "sysadmin",
            "path": "/sysadmin"
        }
    ],
    "scopes": [
        {
            "name": "read.all",
            "displayName": "Allows read all data"
        },
        {
            "name": "write.all",
            "displayName": "Allows writing all data"
        },
        {
            "name": "manage.config",
            "displayName": "Allows management of system configuration"
        }
    ],
    "resources": [
        {
            "name": "data",
            "displayName": "Data access resource",
            "scopes": [
                {
                    "name": "read.all"
                },
                {
                    "name": "write.all"
                }
            ]
        },
        {
            "name": "system-configuration",
            "displayName": "System configuration management",
            "scopes": [
                {
                    "name": "manage.config"
                }
            ]
        }
    ],
    "policies": [
        {
            "name": "is_system_admin",
            "type": "group",
            "description": "Policy detecting if user is system administrator, giving him total control over system configuration",
            "groups": [
                {
                    "path": "/sysadmin",
                    "extendChildren": true
                }
            ]
        },
        {
            "name": "is_data_user",
            "type": "group",
            "description": "Policy detecting if user can read data",
            "groups": [
                {
                    "path": "/data-user",
                    "extendChildren": true
                }
            ]
        },
        {
            "name": "is_data_provider",
            "type": "group",
            "description": "Policy detecting if user can create data",
            "groups": [
                {
                    "path": "/data-user/data-creator",
                    "extendChildren": true
                }
            ]
        }
    ],
    "permissions": [
        {
            "name": "administer-system-permission",
            "description": "Gives control over system configuration",
            "type": "scope",
            "resources": [
                "system-configuration"
            ],
            "scopes": [
                "manage.config"
            ],
            "policies": [
                "is_system_admin"
            ]
        },
        {
            "name": "data-read-permission",
            "description": "Allows to access data in the system",
            "type": "scope",
            "resources": [
                "data"
            ],
            "scopes": [
                "read.all"
            ],
            "policies": [
                "is_data_user"
            ]
        },
        {
            "name": "create-data-permission",
            "description": "Allows to create data in the system",
            "type": "scope",
            "resources": [
                "tasks"
            ],
            "scopes": [
                "write.all"
            ],
            "policies": [
                "is_data_provider"
            ]
        }
    ]
}
```
will generate following json diff:
```json
{
    "client": {
       "op": "ADD",
       "clientSpec": {
          "attributes": {
             "login_theme": "base"
          },
          "authorizationServicesEnabled": true,
          "clientId": "sample-client",
          "enabled": true,
          "protocol": "openid-connect",
          "publicClient": false,
          "redirectUris": [
             "*"
          ],
          "serviceAccountsEnabled": true
       }
    },
    "scopes": [
       {
          "op": "ADD",
          "scopeSpec": {
             "displayName": "Allows read all data",
             "name": "read.all"
          }
       },
       {
          "op": "ADD",
          "scopeSpec": {
             "displayName": "Allows writing all data",
             "name": "write.all"
          }
       },
       {
          "op": "ADD",
          "scopeSpec": {
             "displayName": "Allows management of system configuration",
             "name": "manage.config"
          }
       }
    ],
    "resources": [
       {
          "op": "ADD",
          "resourceSpec": {
             "displayName": "Data access resource",
             "name": "data",
             "scopes": [
                {
                   "name": "read.all"
                },
                {
                   "name": "write.all"
                }
             ]
          }
       },
       {
          "op": "ADD",
          "resourceSpec": {
             "displayName": "System configuration management",
             "name": "system-configuration",
             "scopes": [
                {
                   "name": "manage.config"
                }
             ]
          }
       }
    ],
    "policies": [
       {
          "op": "ADD",
          "policySpec": {
             "description": "Policy detecting if user is system administrator, giving him total control over system configuration",
             "name": "is_system_admin",
             "type": "group",
             "groups": [
                {
                   "path": "/sysadmin",
                   "extendChildren": true
                }
             ]
          }
       },
       {
          "op": "ADD",
          "policySpec": {
             "description": "Policy detecting if user can read data",
             "name": "is_data_user",
             "type": "group",
             "groups": [
                {
                   "path": "/data-user",
                   "extendChildren": true
                }
             ]
          }
       },
       {
          "op": "ADD",
          "policySpec": {
             "description": "Policy detecting if user can create data",
             "name": "is_data_provider",
             "type": "group",
             "groups": [
                {
                   "path": "/data-user/data-creator",
                   "extendChildren": true
                }
             ]
          }
       }
    ],
    "permissions": [
       {
          "op": "ADD",
          "permSpec": {
             "description": "Gives control over system configuration",
             "name": "administer-system-permission",
             "policies": [
                "is_system_admin"
             ],
             "resources": [
                "system-configuration"
             ],
             "scopes": [
                "manage.config"
             ],
             "type": "scope"
          }
       },
       {
          "op": "ADD",
          "permSpec": {
             "description": "Allows to access data in the system",
             "name": "data-read-permission",
             "policies": [
                "is_data_user"
             ],
             "resources": [
                "data"
             ],
             "scopes": [
                "read.all"
             ],
             "type": "scope"
          }
       },
       {
          "op": "ADD",
          "permSpec": {
             "description": "Allows to create data in the system",
             "name": "create-data-permission",
             "policies": [
                "is_data_provider"
             ],
             "resources": [
                "tasks"
             ],
             "scopes": [
                "write.all"
             ],
             "type": "scope"
          }
       }
    ],
    "groups": [
       {
          "op": "ADD",
          "groupSpec": {
             "name": "data-user",
             "path": "/data-user"
          }
       },
       {
          "op": "ADD",
          "groupSpec": {
             "name": "data-creator",
             "path": "/data-user/data-creator"
          }
       },
       {
          "op": "ADD",
          "groupSpec": {
             "name": "sysadmin",
             "path": "/sysadmin"
          }
       }
    ]
 }
```
## Adding scope and new resource to existing client
Let's assume we alrady have client *sample-client* in realm *sample-realm* with configuration applied from  ([this example](examples/data-access-client-conf.json)). We add new entry with scope and resource to it:
```json
{
    ...
    "scopes": [
       {
        "displayName": "New fancy scope",
        "name": "fancy.scope"
       },
    ...
    ],
    "resources": [
        {
            "name": "fancy-new-resource",
            "displayName": "Completly new resource",
            "scopes": [
                {
                    "name": "fancy.scope"
                }
            ]
        },
    ...
    ]
}
```
and execute diff command which will result in following json:
```json
{
   "client": {
      "op": "NONE",
      "clientSpec": {
         "clientId": "sample-client"
      }
   },
   "scopes": [
      {
         "op": "ADD",
         "scopeSpec": {
            "displayName": "New fancy scope",
            "name": "fancy.scope"
         }
      }
   ],
   "resources": [
      {
         "op": "ADD",
         "resourceSpec": {
            "displayName": "Completly new resource",
            "name": "fancy-new-resource",
            "scopes": [
               {
                  "name": "fancy.scope"
               }
            ]
         }
      }
   ]
}
```
## Removing policy and permission
Let's assume we alrady have client *sample-client* in realm *sample-realm* with configuration applied from  ([this example](examples/data-access-client-extension.json)). We remove permission named *administer-system-permission* and execute diff command will result in following json:
```json
{
   "client": {
      "op": "NONE",
      "clientSpec": {
         "clientId": "sample-client"
      }
   },
   "permissions": [
      {
         "op": "DEL",
         "permSpec": {
            "decisionStrategy": "UNANIMOUS",
            "description": "Gives control over system configuration",
            "id": "a1323329-6f5c-477d-967d-89297234d00d",
            "logic": "POSITIVE",
            "name": "administer-system-permission",
            "type": "scope"
         }
      }
   ]
}
```

