# Project's Architecture

First of all, it's absolutely spherical project in vacuum. Project provided as-is and some architectural decisions can
be disputable. I'll try to describe some ideas here.

## Transactions And Repository Pattern

`TBD`

## Mocking

`TBD`

## Scheme Simplicity

Current implementation is relatively simple in terms of data scheme. What can we do (probably) better?

1. Create initial balance as payment record in database. With (e.g.) zero payer.  
It cause consistency between balance column in accounts table and summary of all payments for account.
1. No transactions or different transaction isolation level. It can improve database performance.    
But possible inconsistency, etc.
1. Caching of running balance with the current scheme. E.g. in redis.  
It will decrease access to database. But cache invalidation will be our new headache.
1. Foreign key on IDs. It's not always good to use foreign keys in highload environment. But hey guys! It's the only
simple sample project! Let it be!

## Pagination In GET Requests

Yeap, we need pagination in actual application ) Page via URL parameters and blah-blah-blah.
