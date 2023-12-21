# Korora Design Spec

_Still in progress_

Korora (the Māori name for the Little Penguin) is an ActivityPub/ActivityStream powered social app.

This document is meant to define what the absolute minimum Korora needs to be able to do to officially be considered "beta" software, and ready for people to start running and playing with.
Along the way, discussions of future possibilities will be included.
These possibilities will be clearly labeled as such, and are not meant to be definitive as any sort of roadmap or future feature list.
They're included simply as notes and reference for future design specs.

## High-level overview

The initial version of Korora is intended to be a self-hosted microblogging platform, in the same genre as Mastodon, GoToSocial, or Takahe.

The key differentiator of Korora from other ActivityPub-based microblogging platforms is the addition of Huddles.
Huddles can be though of like Google Plus' Circles feature.

## Huddles
Huddles (named because Korora are a type of penguin, and a group of penguins is a huddle!) are a way to organize your social graph.

Huddles are meant to define how different people are allowed to interact with you, and what posts they're allowed to see.
When you share a new post, you can define which huddles you want to share it with.
Some example huddles might be:

-   "software nerds", a huddle for people you talk about software development with.
    If someone connects with you after a professional event, you might add them to this huddle.

-   "leet gamerz", a huddle for people you play games with, or talk about gaming with.

-   "baby pictures", maybe you want to share pictures of your kids, but only with close friends and family, this huddle would be for that.

In addition to huddles you organize, you can also define "public huddles", which allow people to subscribe to different subsets of your posts.
Public huddles can be set to allow anyone to follow, or require approvals, similar to Mastodon profiles.
Where you might define your private huddles as above, you could define public huddles like:

-   "shitposting", a huddle someone could follow to only see your shitposts.
-   "dog pictures", a huddle that's just a feed of pictures of your dog.
-   "portfolio", a huddle for people who want to see just your artwork, but not your other posts.


### Huddles in ActivityPub

On a technical level, huddles will probably work something like this:

-   `Create` activities will `cc` a huddle's collection, instead of the `followers` collection.
    The `to` field will only be set to `as:Public` if a post is set to `Public` instead of specific huddles.

    -   **Note** in reading the mastodon source, this may not be feasible, and the `cc` field may just need to be every user on a given instance (or if going to not a shared inbox, `cc` to just the user in question). I can't tell if Mastodon actually resolves collections other than the "magic" `followers` collection.

-   For backwards compatibility, public huddles (those which other people can subscribe to themselves, instead of being placed in manually by you) will be implemented as an Actor.
    Specifically, the actor's `@type` will be a `["as:Person", "korora:Huddle"]`, so that existing software should interpret it as a `Person`, but if other software implements the `korora` namespace it will be able to understand the distinction.

    -   A `Follow` activity will be sent to the huddle's inbox, if the huddle allows anyone to follow, an automatic `Accept` will be sent back, if not the owner will need to approve or deny the request.

    -   When a user `Create`s a new activity to a public huddle, the ideal situation would be the actual user Actor sends the new activity to everyone in a given huddle.
        There's a question of if this will look like that user sending a DM to a given mastodon user, if this looks like direct messages Korora will need to instead send an `Announce` activity to the huddle Actor's followers list. Reading the code, it _looks_ like it shouldn't look like a DM or mention at all, though it does mean that the visibility is `Limited`, which may mean it shows up on a user's timeline on the remote instance? To be explored. Okay so it definitely _looks_ like it should only show up in the feeds of those in the `cc` or `to` field, and shouldn't show up on the public timeline. It also seems like `Limited` status posts aren't boostable in Mastodon? But I think that's fine and if this takes off we can add an extension to the activities that specifies if it's sharable. Neat!

    -   The `korora:Huddle` Actor's outbox should just be `Announce` activities of the original posts by the primary user Actor.
        Ideally, this should mean that remove users browsing the huddle Actor should see the posts shared to that huddle. Again, TBD based on what happens in the real world.


## Minimum UI

A recent trend in the ActivityPub world is to implement the Mastodon API as a means of immediate compatibility with many great existing clients.
Unfortunately, because Korora's design makes decisions that are incompatible with the Mastodon API this is not an option for us.
As such, Korora's minimum UI will need to handle the following tasks:

-   Account management
    -   Configure basic account information (bio, display name, etc)
    -   Ability to require approvals for followers
-   Huddle management
    -   Create new Huddles, set huddle visibility, manage users in huddles
    -   Ability to require approvals for publicly visible huddle followers
-   Public profile and huddle views
-   Timeline view for logged in users
-   Create post UI (with support for images and content warnings)

As an aside, I'm intrigued by the idea of making Korora's backend ActivityPub only, and implementing the frontend as ActivityPub C2S, but I'm concerned that may be an additional level of effort on top of all of the rest of the app.
