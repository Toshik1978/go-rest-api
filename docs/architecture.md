# Project's Architecture

This document is a kind of FAQ by project's architecture.

## Restrictions On Server Side

_Why did you use columns restrictions in database instead of programming logic?_

Obviously, programming logic gives more flexibility. Especially in error handling.
But any programming restriction either not consistent, or use some locking. Database already consistent solution.
So, why not to use for simplicity.

## Transactions And Repository Pattern

_Don't you think, that transactions badly lives with Repository Pattern?_

Yes, I think so. That's why we should not use some explicit transaction objects.
We are using kind of unit-of-work, but probably smth. like Command Pattern can better fit requirements.
Anyway, Repository Pattern is extremely simple. And unit-of-work runs transparently on the top of it.

## Scheme Simplicity

_You are always talking about simplicity. Why do you insist on this?_

Because simple logic produce less support time. Look at some changes to data scheme, for example.

1. Create initial balance as payment record in database. With (e.g.) zero payer.
It cause consistency between balance column in accounts table and summary of all payments for account.  
This way we should have some special kind of payment (e.g. with zero payer). Nullable foreign key, etc.
What is more important it produce additional fake payment in payments list. Do we need to filter it?
1. Don't use transactions or modify transaction isolation level. It can improve database performance/simplify logic.  
Eventually consistency in the better case. Data inconsistency in the worst case.
1. Caching of running balance with the current scheme. E.g. in redis.  
It will decrease access to database. But cache invalidation will be our new headache.
1. Don't use foreign keys on IDs. It's not always good to use foreign keys in highload environment.  
Possible data inconsistency. Required to do more checks manually. E.g. checking of account existence during payment.
This way we will have more database queries and more ways to do mistake w/o locking (add payment in parallel with account deletion).
1. Add pagination in GET requests.  
Yes, we need pagination in actual application. Agree )

## Kind Of Dependency Injection

_Looks like main function initializes most of dependencies and pass it forward. Why do you rely on non-nil initialization?_

Probably we can panics everywhere to check, that `New...` takes non-nil input. It will require more logic in unit tests.
Now we can pass `nil` in tests, not requiring specific field to be initialized. With panic we should more objects everywhere.

Besides I'm prefer to use functional tests covers all initialization logic.

Anyway, it's a point to discuss with reviewer )