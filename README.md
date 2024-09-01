# Bythen Take Home Test


Halo Kak Gemma and Kak Ridhwan,

I want to say sorry, because i can't completing the task for now,  detail of problem already share on email.

but i'm still trying to complete the task even i can't continue to next step, because i want to research about backend template for my project or next project at internal or external team.

best regards,
Ananda Affan Fattahila

---

This is link to access api spec documentation
https://api.bythen.fanzru.dev/doc/swagger

---
## Migration DB Documentation

database support mysql 8.0

> If you have change in database please create sql syntax to rollback, because infra need sleep.

Please Install  `golang-migrate`  to migrate database, to install  `golang-migrate`  you can read this documentation:

**Windows :**

```
https://verinumbenator.medium.com/installing-golang-migrate-on-windows-b4b3df9b97b2
```

**Mac :**

```
brew install golang-migrate
```

**Ubuntu :**

```
https://www.geeksforgeeks.org/how-to-install-golang-migrate-on-ubuntu/
```

## Migration Script


if have done to install  `golang-migrate`, please prepare your database and create new database, for example  `kerjago_db`.

**Migrate Create**

Create a migration file. You can find the file at  `migration`  folder

```
make migrate-create NAME=namefile
```

**Migrate Up**

To migrate all your migration file

```
make migrate-up
```

**Migrate Down**

To delete all your schema with migration

```
make migrate-down
```

**Migrate Rollback**

to run migration down only  `N`  step(s) behind

```
make migrate-rollback N=yournumberrunmigrationdown
```

**Fixing your Migration**

What happend if your database is dirty?

You can fix your migration first and then using foce command with the version you want.

If you're happend to get  `error: Dirty database version 16. Fix and force version.`

Then you want to run:

```
make migrate-force VERSION=1
```

Reference:  [golang-migrate/migrate#282 (comment)](https://github.com/golang-migrate/migrate/issues/282#issuecomment-530743258) 