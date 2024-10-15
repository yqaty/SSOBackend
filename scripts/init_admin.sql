INSERT INTO public."user" (created_at, updated_at, uid, email, name)
VALUES (
        current_timestamp,
        current_timestamp,
        'ffb6e834-3615-4ebb-9d9d-825af333a3ca',
        'root@hustunique.com',
        '联创团队'
    );
INSERT INTO public.user_role (created_at, updated_at, role_name, uid)
VALUES(
        current_timestamp,
        current_timestamp,
        'admin',
        'ffb6e834-3615-4ebb-9d9d-825af333a3ca'
    );
INSERT INTO public.object (created_at, updated_at, action, resource)
VALUES (
        current_timestamp,
        current_timestamp,
        '.*',
        'UniqueSSO::Internal::.*'
    );
INSERT INTO public.object_group (
        created_at,
        updated_at,
        object_group_name,
        object_id
    )
SELECT current_timestamp,
    current_timestamp,
    'UniqueSSOFullInternalAPI',
    id
from public.object
WHERE action = '.*'
    AND resource = 'UniqueSSO::Internal::.*';
INSERT INTO public.permissions (
        created_at,
        updated_at,
        role_name,
        object_group_name
    )
VALUES (
        current_timestamp,
        current_timestamp,
        'admin',
        'UniqueSSOFullInternalAPI'
    );