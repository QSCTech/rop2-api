SET @orgName = "要新建的组织名";
SET @departName = "要新建的部门名";
SET @adminId = "admin_id";
SET @adminNickname = "初始管理员";

INSERT INTO orgs (`name`,create_at) VALUES (@orgName,now());
SELECT LAST_INSERT_ID() INTO @orgId;

-- 新建默认部门
INSERT INTO departs (`name`,create_at,`owner`) VALUES (concat(@orgName,"-默认部门"),now(),@orgId);
SELECT LAST_INSERT_ID() INTO @defaultDepartId;
UPDATE SET default_depart = @defaultDepartId WHERE id = @orgId;

-- 新建普通部门
INSERT INTO departs (`name`,create_at,`owner`) VALUES (@departName,now(),@orgId);

-- 新建管理员(权限级别20)
INSERT INTO admins (`zju_id`,`at`,nickname,`level`,create_at) VALUES (@adminId,@orgId,@adminNickname,20,now());