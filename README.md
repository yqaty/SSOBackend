# UniqueSSO - Backend

UniqueStudio Single Sign On system and uniform permission control system.

## Single Sign On

The SSO feature is designed for [Traefik ForwardAuth middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) while Traefik is the gateway front of all services provided by UniqueStudio. **The standard SSO implementation will be supported in next major version.**

The supported sign in methods are shown below:

- mobile & password & *2FA(optional)*
- email  & password & *2FA(optional)*
- mobile & SMS
- email  & validation code
- Oauth
  - lark

## Permission Control

The permission control system is a basic implementation of RBAC. A specification is shown below:

- User
    The subject that needs some permissions to do some things
- Role
    An entity bound a group of permissions
- Object
    The combination of `Action` and `Resource`
- Object Group
    a bunch of objects

In semantic level, a permission is a triple: `<role, action, resource>`, which means the specific `role` contains ability to perform the `action` to `resource`.

For flexibility, an user can bind multiple roles while a role can be granted with multiple object groups. In addition, to shrink duplicate object creations, a role can be only bound to object groups instead of objects directly.

For example, user `xylonx` is granted as role `admin` while `admin` is granted with object group `OpenPlatform` with object `<Push, OpenPlatform::SMS>` and object group `UniqueSSO::Internal` with objects `<Grant, UniqueSSO::Internal::User::Role>`, `<Revoke, UniqueSSO::Internal::User::Role>`. Therefore, `xylonx` can push sms, assign and revoke roles to someone.



<u></u>
