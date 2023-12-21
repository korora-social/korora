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

### Huddles and Announces/Boosts/Reblogs/Reposts

Because boosts/reblogs/reposts/announces are a user expectation, and a feature I'd like to support, the Korora timeline needs some mechanism to easily enable them without a cumbersome user experience.

Similar to the Twitter and Bluesky systems, hitting the boost icon will display a modal with several options.
The options should be:

-   "Boost or Reblog...", which opens a compose editor allowing you to specify which huddles to send the boost to, or add an additional, Tumblr-style reblog.
-   "Boost to $HUDDLE_NAME", an option that should be repeated for each huddle the person you're boosting it from (or the person who boosted it to you) is in.

Users should additionally have the ability to disable the "quick-boost" functionality, and just always be taken into the boost/reblog compose screen.

### Organizing new relationships into Huddles

When you follow someone, or they follow you, that relationship should be immediately placed into an "Unorganized Relationships" huddle.
At a later point, you should be able to review your unorganized relationships, with a prompt to organize them into your huddles.

The Korora UI should immediately prompt you to categorize users into huddles when you click the "Follow" button (or perhaps this should be replaced with "Add to huddles..."/"Manage Huddles").

If your account is set to require approvals for follow requests, the approvals screen should prompt you to categorize newly approved followers into huddles.

### Public Huddles

You should be able to declare certain huddles as public, meaning other people are able to follow or unfollow that huddle at will.
However, you should still be able to organize people who follow your user into those huddles, so that you don't need to have a second, private, huddle for the same purpose.

Internally, Korora should note if someone is a follower of the public huddle, or a follower of your main user that you've also put in that huddle.
Applications which understand the huddle paradigm should allow people who've been added to the "private" side of the huddle to remove themselves from that huddle, if they want.
On a technical level, this should be achieved by the huddle's `Followers` collection being accessible (with security considerations, may just show subsets when requested) to someone who has been added to that huddle.
To be removed from a public huddle, just send an Unfollow to that huddle's inbox.

### Huddles in ActivityPub

On a technical level, huddles will probably work something like this:

-   `Create` activities will `cc` a huddle's collection, instead of the `followers` collection.
    The `to` field will only be set to `as:Public` if a post is set to `Public` instead of specific huddles.

    -   **Note** in reading the mastodon source, this may not be feasible, and the `cc` field may just need to be every user on a given instance (or if going to not a shared inbox, `cc` to just the user in question). I can't tell if Mastodon actually resolves collections other than the "magic" `followers` collection.

    -   Posts sent to public huddles will have a `cc` of that huddle's `followers` list

-   For backwards compatibility, public huddles (those which other people can subscribe to themselves, instead of being placed in manually by you) will be implemented as an Actor.
    Specifically, the actor's `@type` will be a `["as:Person", "korora:Huddle"]`, so that existing software should interpret it as a `Person`, but if other software implements the `korora` namespace it will be able to understand the distinction.

    -   A `Follow` activity will be sent to the huddle's inbox, if the huddle allows anyone to follow, an automatic `Accept` will be sent back, if not the owner will need to approve or deny the request.

    -   When a user `Create`s a new activity to a public huddle, the ideal situation would be the actual user Actor sends the new activity to everyone in a given huddle.

        There's a question of if this will look like that user sending a DM to a given mastodon user, if this looks like direct messages Korora will need to instead send an `Announce` activity to the huddle Actor's followers list.
        Reading the code, it _looks_ like it shouldn't look like a DM or mention at all, though it does mean that the visibility is `Limited`, which may mean it shows up on a user's timeline on the remote instance? To be explored.
        Okay so it definitely _looks_ like it should only show up in the feeds of those in the `cc` or `to` field, and shouldn't show up on the public timeline.
        It also seems like `Limited` status posts aren't boostable in Mastodon? But I think that's fine and if this takes off we can add an extension to the activities that specifies if it's sharable.
        Neat!

    -   The `korora:Huddle` Actor's outbox should just be `Announce` activities of the original posts by the primary user Actor.
        Ideally, this should mean that remote users browsing the huddle Actor should see the posts shared to that huddle. Again, TBD based on what happens in the real world.

-   A new `korora:boostable` attribute will be added to post activities that are meant to be displayed in a timeline.
    This attribute will specify if a post is meant to be shared by others.
    The reason for making this an explicit attribute, instead of deriving it from other attributes into Mastodon's `visibility` attribute is to allow for other software to better interpret a post author's intent.
    Because most Korora posts will be defined by Mastodon as `limited` visibility, Mastodon users won't be able to boost posts shared to different circles.
    This is less than ideal because I may want to share something just to one huddle, but allow members of that huddle to further share the post.

-   Korora will still need to keep track of the traditional Followers/Following paradigm, because while you may be following someone, they may not follow you.
    This is because you may send a post to a huddle which includes people not following you or that huddle.
    In this case, Korora may include them in `to` and `cc` for outboxes, but it should not deliver those users' inboxes.


## Minimum UI

A recent trend in the ActivityPub world is to implement the Mastodon API as a means of immediate compatibility with many great existing clients.
Unfortunately, because Korora's design makes decisions that are incompatible with the Mastodon API this is not an option for us.
As such, Korora's minimum UI will need to handle the following tasks:

-   Account management
    -   Configure basic account information (bio, display name, etc)
    -   Ability to require approvals for followers
-   Basic User Safety
    -   Ability to run in a limited federation mode (allowlist instances)
    -   Ability to block instances
    -   Ability to block individual users (possibly a "Blocked" huddle?)
    -   Ability to "softblock" users (remove them as a follower, but not fully block them; possibly just block + unblock actions)
    -   Ability to limit DMs/mentions to everyone, no one, followers only, specific huddles, specific users
        -   Similarly, an inverse ability, to block DMs/mentions from specific huddles or specific users
-   Huddle management
    -   Create new Huddles, set huddle visibility, manage users in huddles
    -   Ability to require approvals for publicly visible huddle followers
-   Public profile and huddle views
-   Timeline view for logged in users
    -   Timeline of everything
    -   Timelines for specific huddles
-   Create post UI (with support for images and content warnings)

The frontend should communicate with the backend entirely through a public API that other applications should be free to consume.

As an aside, I'm intrigued by the idea of making Korora's backend ActivityPub only, and implementing the frontend as ActivityPub C2S, but I'm concerned that may be an additional level of effort on top of all of the rest of the app.


## Mastodon Federation Compatibility

Because Mastodon is the... mastodon of ActivityPub-based microblogging platforms, a certain level of federation compatibility is a desired end goal.

Korora's explicit interoperability goals with Mastodon are:

-   Posts sent to huddles are visible in members of that huddle's timeline
-   Replies to posts in a huddle are properly federated outward to members of that huddle
-   Korora users can follow and interact with Mastodon users
-   Korora users who reply to a post from a Mastodon user can limit the reply's visibility to everyone/author only/author's followers/author's followers + specific huddles, and that limit works as intended
-   Mastodon users should be able to subscribe to public huddles, and see posts sent to those huddles in their timelines
-   Webfinger `acct:` support, even if you really shouldn't identify users by a webfinger `acct:` url.

Korora's explicit non-goals for interoperability with Mastodon are:

-   Korora will not support the Mastodon API, or even a subset of it.
-   Because of Korora's ideas around more limited access, Korora does not intend to figure out a hack to allow Mastodon users to boost posts sent to private huddles.


## Outstanding questions

_Various design decisions that need to be made, but I don't have an answer for yet_.

-   If you remove someone from all huddles, should that be considered an unfollow and softblock?
    -   I'm leaning towards no, and they'll just go back into an "Unorganized" huddle, and should still be able to see things posted to "everyone" or "public".
        This would make huddles more of an opt-in feature, so you could still uses Korora in the same way as other microblogging platforms.
        It would also allow you to dip your toes into using huddles, start by creating just one huddle to use like Instagram's close friends, and then maybe start adding more as you find useful.
