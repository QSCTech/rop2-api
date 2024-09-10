-- 查询全数据库内，提交过表单的学号数，和未提交过表单(但已登录过)的学号数
-- 拥有任意管理权限的学号不参与上述统计

select 
  count(case when distinct_results.create_at is not null and admins.nickname is null then 1 end) as submitCount,
  count(case when distinct_results.create_at is null and admins.nickname is null then 1 end) as restCount,
  concat(substring(people.zju_id,1,3), repeat('x',length(people.zju_id)-3)) as pattern from people
left join (
  select distinct zju_id, create_at from results
)`distinct_results` on distinct_results.zju_id = people.zju_id 
-- 管理员无需去重
left join admins on admins.zju_id = people.zju_id
group by pattern;