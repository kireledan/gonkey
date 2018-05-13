# Gonkey

I don't like the other configurators out there. They are either too complex for what they need to do, or are slow. Gonkey is meant to be a go-learning project and a desire to create a configurator that's light, simple, and easy to extend. 

## To Run

To use it, simply run. Here's an example yaml file.

    gonkey <my.yaml>



    tasks:
      - name: Creating file derp
        module: command
        args:
          cmd: "touch derp"
      - name: Copying derp
        module: copy
        args:
          src: "./derp"
          dst: "./derpcopy"
    



## General Guidelines

All tasks are declared inside a tasks list. The program executes the corresponding module from the module name and runs with the given args. These modules are not case sensitive.

The currently available modules are:

- Command - run a command
- lineinfile - replace a line in a file with given regex
- Service (experimental) - set the desired status of a service
- template - copy a template file with the proper values
- Copy - copy a file with desired mode,owner,group

#### Command

the command module only takes 1 argument, cmd. This is the command that is executed

    tasks:
      - name: Creating file derp
        module: command
        args:
          cmd: "touch derp"



#### lineinfile

This command will replace the entire line of that matches the search pattern.

    tasks:
      - name: Replace the <changeme> with changed
        module: lineinfile
        args:
          file: "tobechanged.txt"
          search: "<changeme>"
          replacewith: "changed"

#### Service

This will run a particular service command. Currently this only supports systemd and homebrew services.

    tasks:
      - name: Start the etcd service
        module: service
        args:
          name: etcd
          status: start



#### Template

This writes a template to a desired directory given a file with values.

    tasks:
      - name: Writing values to my template
        module: template
        args:
          template: "./mytemplate.temp.conf"
          output: "./mytemplate.conf"
          valuefile: "./myvars.yaml"

The value file looks like, 

    # ./myvars.yaml
    args:
      woo: "it worked!"
    

The template looks like

    # ./mytemplate.temp.conf
    this is atemplate
    a value goes here ({{.GetConfig "woo"}})
    

All variables from a valuefile are grabbed via .GetConfig "<desiredvalue>". Any valid GoTemplate syntax will work for these templates too. 



#### Copy

This will copy a file from src to dest. 

    tasks:
      - name: Copying derp
        module: copy
        args:
          src: "./derp"
          dst: "./derpcopy"

If you want to change the mode,user, or group, just add the arguments

    tasks:
      - name: Copying derp
        module: copy
        args:
          src: "./derp"
          dst: "./derpcopy"
          mode:
          owner:
          group:





## Hacking.

To test

    dep ensure
    go test ./...


