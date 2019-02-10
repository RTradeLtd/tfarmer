# tfarmer

In order to protect user privacy and data, we avoid integration with any third-party data collection tools like google analytics. The unfortunate side effect to this is that determining how active Temporal is (amount of data going through our systems, daily active user counts) is a lot harder, and not as smooth. Using `tfarmer` we can scrape our database to determine activity, without having to expose the data to third party services, and risk compromising the privacy of our users.

Some notes about what `tfarmer` will and won't do

## do's and donts

* daily active users:
  * do:
    * determine daily active users based on when they last signed in
    * count the number of users to determine the number
    * email the number of users to the CTO+CEO
  * dont:
    * run outside of our secure data center environment
    * send any information other than the **number** of users (this means we wont be sending information like usernames, emails, etc..)

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
