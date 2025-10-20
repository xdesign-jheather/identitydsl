# AWS IAM Identity Center Accelerator

`identitydsl` is a free utility to convert a concise human-readable configuration for AWS IAM Identity Center into machine readable IaC, namely terraform.

## Why?

Managing users, groups, accounts and assignments for medium to large size businesses can become complex and repetitive at scale, with a high proliferation of IaC resources. `identitydsl` allows a more concise configuration (using the DSL) and helps you understand who has access to what.

The focus is on ease of human authorship and review, which happens to convert to dependable, consistent and optimized IaC.

## IaC

To start with we just support one IaC solution: [terraform](https://www.hashicorp.com/en/products/terraform). The output format is [JSON](https://developer.hashicorp.com/terraform/language/syntax/json) based terraform, but [HCL](https://developer.hashicorp.com/terraform/language/syntax/configuration) will be supported in the future.


## The DSL

No surprises that `identitydsl` defines a small and simple text based configuration language.

The primary aim of the DSL is to describe key entities such as Users, Groups and Accounts, then decorate them with custom tags and/or labels. These tags and labels can be used to make batch assignments by filtering.

### Comments

You can add comments:

```
// A comment line starts with two slashes
```

### Accounts

To create an account:

```
Account 12345678
```

To add labels to an account:

```
Account 12345678
   Label1
   Label2
   "Label 3"
```

To add tag to an account:

```
Account 12345678
   Name Sales
   Contact "Bobby Tables"
```

To add labels or tags to multiple accounts at once:

```
Account 12345678, 87654321
   Owner Legal
   Environment production
```

### Users

To create a simple user:

```
User Jonathan
```

Tags and labels can be added to one or more groups as described with `Account`:

```
// Add the Team tag and BillingAccess label to 3 users at once

User A, B, C
  Team Platform
  BillingAccess
```

### Groups

To create a simple group:

```
Group Developers
```

Tags and labels can be added to one or more groups as described with `Account`:

```
// Add the Department tag and BillingAccess label to 3 groups at once

Group FooDevelopers, FooTesters, FooFighters
  Department Foo
  BillingAccess
```

### Roles

A role represents the permission set in Identity Center.



To create a simple role:


```
// This will create a `ReadOnly` permission set with a `ReadOnly` policy.

Role ReadOnly
```

To add additional/alternative policies to a role:

```
// This will create a `ReadOnly` permission set with `OtherPolicy` and `AnotherPolicy` policies attached.

Role ReadOnly
  OtherPolicy
  AnotherPolicy
```

Short names indicate these are policies you would like to implement using this IaC output.

If you want to use an existing customer managed or AWS managed policy instead, just specify the ARN:

```
Role FullAccess
  arn:aws:iam::aws:policy/AmazonEC2FullAccess
```

> **_Note_** Roles cannot be grouped with tags or labels. This is intentional and each assignment must be explicit and intentional.

### Labels

Labels are keyless string associations for grouping entities. We support spaces in the label using double quotes around the whole label value.

```
// Account `12345` shows 2 ways to define labels

Account 12345
    RDS
    "Website DR"
```

> **_Note_** Labels that match IDs are not permitted and will produce an error.


### Tags

Tags are key value pairs used for filtering through entities, much like AWS. We support spaces in the key and the value using double quotes around both.

```
// User `Bob` shows 4 ways to define tags

User Bob.Smith
    Email bob@example.com
    "Full Name" "Bob Smith"
    "Given Name" Bob
    Access "Developer Level 4"
```

> **_Note_** In practice we expect most organisations use one or the other, but we support both. 

### Assignments

Everything stated above is designed to make assignments work effectively.

Each assignment must select:

- One or more `Account`
- One or more `User` or `Group`
- One or more `Role`

Selections can be expressed with:

- Name or ID
- A label
- A tag key value pair

This produces a set of assignments for Identity Center that is the [cartesian product](https://en.wikipedia.org/wiki/Cartesian_product) of:

`Account * (User + Group) * Role`

To express a single assignment:

```
// This will produce 1 assignment

Assign
  Account Account1
  Role Role1
  Group Group1
```

To express multiple entities, just provide a comma delimited list:

```
// This will produce 8 assignments

Assign
  Account Account1, Account2
  Role Role1, Role2
  User User1, User2
```

Rather than listing these explicitly, you can make use of your tags and labels for `Account` and `Group`:

```
// Tag two accounts with Owner

Account 1, 2
  Owner Legal
  
// Assign to accounts with Owner = Legal (1 and 2 in this case)   
  
Assign
  Account Owner Legal
  Role Role1
  Group Group1
```

You can go one step further and specify multiple filters, consisting of tags, labels or IDs:

```
// Some accounts

Account Hello, World
  Team Data
  
Account Foo
  Snowflake

Account Bar
  
// Assign to 4 accounts using 3 methods of selection:
  
Assign
  Account Team Data, Snowflake, Bar
  Role DBAReadOnly
  Group DBA
```

### Contexts

A context is a way of expressing multiple similar assignments, without repeating `Account`, `User` or `Group` selections.

For example, the following two examples are functionally the same:

```
Assign
  Account Team Data
  Role ReadOnly
  Group DataTeam

Assign
  Account Team Data
  Role ReadOnlyGuest
  Group AnotherTeam
```

```
Accounts Team Data

  Assign
    Role ReadOnly
    Group DataTeam

  Assign
    Role ReadOnlyGuest
    Group AnotherTeam
```

You can still specify `Account` in the assignment if you wish, but the context will apply it's filter first.

This is a handy way to vary access in different accounts for example:

```
Accounts Team Data

  Assign
    Account Environment Dev
    Role ReadWrite
    Group DataTeamOperations, DataTeamDeveloper

  Assign
    Account Environment Production
    Role ReadOnly
    Group DataTeamDeveloper

  Assign
    Account Environment Production
    Role ReadWrite
    Group DataTeamOperations
    
```

The example above can be refined further by nesting contexts. The following is functionally equivalent:

```
Accounts Team Data

  Accounts Environment Dev
    
    Assign
      Role ReadWrite
      Group DataTeamDeveloper, DataTeamOperations

  Accounts Environment Production

    Assign
      Role ReadOnly
      Group DataTeamDeveloper

    Assign
      Role ReadWrite
      Group DataTeamOperations
```

## Commands

### validate

The `validate` command will parse your DSL file and check for syntax and logical errors.

```
identitydsl validate ic.txt
```

### synth

The `synth` command will synthesize your DSL and produce IaC output in the working directory in the desired format (only terraform json for now)

```
identitydsl synth ic.txt [-provider=terraform] [-format=json]
```
