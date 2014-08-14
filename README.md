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

* Install <a href="http://golang.org">Golang</a> then add golang binary into system PATH and verify the sucessfull installation by 
<pre>
~$ go version
go version go1.3.1 linux/amd64
</pre>

* Install <a href="http://dartlang.org">Dartlang(Dart SDK)</a> then add dart-sdk binary into system PATH and verify the sucessfull installation by 
<pre>
~$ dart --version
Dart VM version: 1.5.8 (Tue Jul 29 07:05:41 2014) on "linux_x64"
</pre>

* Install <a href="http://www.mysql.com">Mysql</a> then add mysql binary into system PATH and verify the sucessfull installation by 
<pre>
~$ mysql --version
mysql  Ver 14.14 Distrib 5.5.37, for debian-linux-gnu (x86_64) using readline 6.2
</pre>

* Install <a href="http://hbase.apache.org">Hbase</a> then add hbase binary into system PATH and verify the installation is sucessfull by 
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

Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
