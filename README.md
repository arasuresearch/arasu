<h3>Note: arasu is in beta ,still not yet ready for deployment.</h3>


Arasu   :  A Lightning Fast Web Framework
=====
<h4><strong> Note: </strong>arasu development only work on dart enabled browsers like dartium or dart enabled chrome browser.</h4>

Arasu is a Next Generation Full Stack Web framework written on Go language & Dart language.  

Features
========
* lightning fast, because of golang and dartlang
* use RDBMS(mysql is default) and BIGDATA(hbase is default) for serverside store
* use IndexedDB and Angular Dart for clientside store,clientside framework 
* use TDD defaultly supported by golang and dartlang 
* use BDD with selenium and Spinach (this is not yet started github.com/arasuresearch/arasu/bdd)
* automatic build system.

Installation
============

* Install <a href="http://golang.org">Golang</a> then add golang binary into system PATH.

<pre>
export GOROOT="go-installation-directory" 
</pre>
for example if you extracted downloaded go.tag.gz into "/home/user/go" then 
<pre>
export GOROOT=/home/user/go
</pre>
<pre>
export GOPATH="any-directory-where-you-want-to-keep-go-libraries" 
</pre>
for example you can use
<pre>
export GOPATH=/home/user/gopath
</pre>
add both bin into system PATH by
<pre>
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
</pre>

verify the sucessfull installation by 
<pre>
~$ go version
go version go1.4 linux/amd64
</pre>

* Install <a href="http://dartlang.org">Dartlang(Dart SDK)</a> then add DART_HOME & dart-sdk binary into system PATH.

for example if you extracted downloaded dart.tag.gz into "/home/user/dart" then add
<pre>
export DART_HOME=/home/user/dart
export PATH=$DART_HOME/dart-sdk/bin:$PATH
</pre>

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
<pre>
export HBASE_HOME=/hbase/extracted/folder
export JAVA_HOME=/java/extracted_or_installed/folder
</pre>
export PATH=$HBASE_HOME/bin:$JAVA_HOME/bin:$PATH


verify the sucessfull installation  by 
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
</pre>
// you will get output like "You don't have a lockfile, so we need to generate that:" by 
// Dart Pub Manager ,this will take few more seconds (this will occur at first time only).
  
after successfull start....

now visit "http://localhost:4000/" or "http://localhost:4000/#/admins" on Dartium
<pre>
~$ $DART_HOME/chromium/chrome --user-data-dir=$HOME/.config/google-dart http://localhost:4000/#/admins
</pre>

open dartium dev tools settings and disable cache by checking "Disable cache (while DevTools is open)" 



for ubuntu users (tested only on ubuntu 14.10) 
if you get error like
<pre>
dart/chromium/chrome: error while loading shared libraries: libudev.so.0: cannot open shared object file: No such file or directory
</pre>
then check "/lib/x86_64-linux-gnu" folder for latest libudev and link it to 
<pre>
sudo ln -s /lib/x86_64-linux-gnu/libudev.so.1.4.0 /usr/lib/libudev.so.0
</pre>


There you can play !!!



Part 2 
------
Note : for normal application you can use any RDBMS like mysql as we did in the above part 1.
this part 2 is using BigData(hbase) as database.
please use hbase with caution (read more at <a href="http://hbase.apache.org">Hbase</a>)

Creating scaffold for BigData (hbase)

stop the arasu server by pressing CTRL + C

open another new terminal and start bigdata...
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

<pre>
~$ $DART_HOME/chromium/chrome --user-data-dir=$HOME/.config/google-dart http://localhost:4000/#/users
</pre>

there you are !!!!!!!! 
enjoyyyy!!!!!!


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
