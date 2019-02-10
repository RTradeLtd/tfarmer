# tfarmer

`tfarmer` is a data analytics tool used to collect activity/usage information about Temporal, to better allow RTrade Technologies to develop the service, and provide better functionality, support, and usability for frequently used aspects of Temporal.

In order to maintain transparency into our operations, and Temporal as a whole, we believe it is neccessary to also open-source the code we use to gain insight into Temporal while not having to use third-party analytics tooling.

## Overview

In order to protect user privacy and data, we avoid integration with any third-party data collection tools like google analytics. The unfortunate side effect to this is that determining how active Temporal is (amount of data going through our systems, daily active user counts, most frequently used services, etc..) is a lot harder, and can be quite cumbersome. 

Using `tfarmer` we can scrape our database to determine activity, without having to expose the data to third party services, and risk compromising the privacy of our users.

Some notes about what `tfarmer` will and won't do.

tl;dr: collect non-sensitive/identifying information, and email CTO+CEO on a daily basis while conducting all information gathering in our secure data center environment.

## Do's and Dont's

* daily active users:
  * do:
    * determine daily active users based on when they last signed in
    * count the number of users to determine the number
    * email the number of users to the CTO+CEO
  * dont:
    * run outside of our secure data center environment
    * send any information other than the **number** of users (this means we wont be sending information like usernames, emails, etc..)  
<br />
* data flowing through the system:
  * do:
    * determine the most popular data injection types (file uploads, direct pinning)
    * determine how many IPNS records are published each day
    * determine how many pubsub messages are sent each day
    * determine how many keys are created each day
    * determine amount of data stored by file uplods
    * determine amount of data stored by direct pinning
    * email the numbers to the CTO+CEO
  * dont:
    * run outside of our secure data center environment
    * send any information other than the **numbers** (this means we wont be sending information like specific content hashes, specific ipns records, usernames, emails, etc...)
