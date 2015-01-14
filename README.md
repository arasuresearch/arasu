<h3>Note: arasu is in beta ,still not yet ready for deployment.</h3>


Arasu   :  A Lightning Fast Web Framework
=====
<h4><strong> Note: </strong>arasu development only work on dart enabled browsers like dartium or dart enabled chrome browser.</h4>

Arasu is a Next Generation Full Stack Web framework written on Go language & Dart language.  

Features
========
* lightning fast, because of golang and dartlang
* use RDBMS and BIGDATA for serverside store
* use IndexedDB and Angular Dart for clientside store,clientside framework 
* use TDD default by golang and dartlang 
* use BDD with selenium and Spinach (this is in alpha)
* automatic build system.

Installation
============

* Install <a href="http://golang.org">Golang</a> then add golang binary into system PATH.

export GOROOT="go-installation-directory" 
for example if you extracted downloaded go.tag.gz into "/home/user/go" then 
export GOROOT=/home/user/go

export GOROOT="any-directory-where-you-want-to-keep-go-libraries" 
for example you can use
export GOPATH=/home/user/gopath

add both bin into system PATH by
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH

verify the sucessfull installation by 
<pre>
~$ go version
go version go1.4 linux/amd64
</pre>

* Install <a href="http://dartlang.org">Dartlang(Dart SDK)</a> then add dart-sdk binary into system PATH.

add dart bin into system PATH by
for example if you extracted downloaded dart.tag.gz into "/home/user/dart" then 

export PATH=/home/user/dart/dart-sdk/bin:$PATH

verify the sucessfull installation by 
<pre>
~$ dart --version
Dart VM version: 1.8.3 (Mon Dec  1 08:42:49 2014) on "linux_x64"
</pre>

* Install <a href="http://www.mysql.com">Mysql</a> then add mysql binary into system PATH and verify the sucessfull installation by 
<pre>
~$ mysql --version
mysql  Ver 14.14 Distrib 5.5.37, for debian-linux-gnu (x86_64) using readline 6.2
</pre>


* Install <a href="http://hbase.apache.org">Hbase</a> then add hbase binary into system PATH.

set HBASE_HOME , JAVA_HOME and add these bin into system $PATH 
export HBASE_HOME=/hbase/extracted/folder
export JAVA_HOME=/java/installed/folder

export PATH=$HBASE_HOME/bin:$JAVA_HOME/bin:$PATH


verify the installation is sucessfull by 
<pre>
~$ hbase version
//some valid output
</pre>

finally 
* Install Arasu Framework by

> ~$ go get github.com/arasuresearch/arasu 

Creating a New Arasu Project
============================

Part 1  
------
Creating scaffold for relational Database Management System aka RDBMS (mysql)

<pre>

~$ arasu new demo
~$ cd demo
~$ arasu dstore create
~$ arasu generate scaffold Admin name password:string age:integer dob:timestamp sex:bool
~$ arasu dstore migrate  
</pre>

Now start the server:

<pre>

~$ arasu serve

// you will get output like "You don't have a lockfile, so we need to generate that:" by DArt Pub Manager ,this will take few more seconds (this will occur at first time only).
// 
// then 
//you may get dart-sdk "pub downlad error" for few times , but you can ignore and stop the command by CTRL + C .
//and start the same command again until sucessfull start.
  
~$ arasu serve
</pre>
  
after successfull start....

now visit http://localhost:4000/ on 
<i><a href="https://www.dartlang.org/tools/dartium">Dartium</a> or dart enabled chrome</i> 
browser. 
<pre>
To open dartium 
~$ ./DART-SDK-INSTALLED-DIRECTORY/chromiun/chrome 
</pre>

then visit 
> http://localhost:4000/admins.html

There you can play !!!

Part 2 
------
Creating scaffold for BigData (hbase)

stop the arasu server by pressing CTRL + C

open another terminal and start bigdata...
<pre>
~$ start-hbase.sh
~$ hbase thrift start
</pre>

leave this terminal to run thrift deamon. 
come back to old terminal then 

<pre>
~$ arasu dstore create --dstore bigdata
</pre>

this will result in failure

unfortunately Hbase thrift V1 Binary server is not supporting to create database through API Calls
so we have to create it manually . to do that
<pre>

~$ hbase shell
> create_namespace 'demo_development'
> quit

</pre>

close hbase shell , then 

<pre>

~$ arasu generate scaffold User Profile:{FirstName:String,LastName:String,Age:int,Dob:DateTime} Contact:{Phone:String,Email:String} --dstore bigdata
~$ arasu dstore migrate --dstore bigdata
</pre>

Now start the server:
<pre>
~$ arasu serve
</pre>

now visit 
> http://localhost:4000/users.html

on dartium browser, there you can play !!!!! 

<p>lets dive into Full Tutorial  <a href="http://arasuframework.org">Arasu Framework</a> To learn more...</p>


Contribute 
============================
<p>Contribution are welcome <a href="http://arasuframework.org/#contribute">Here</a>.</p>   


License
============================

The BSD License (BSD)

Copyright (c) 2015 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research Lab Pvt Ltd. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above copyright notice, this list of
   conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
   * Neither Thaniyarasu Kannusamy <thaniyarasu@gmail.com>. nor ArasuResearch Inc may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND AUTHOR
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
