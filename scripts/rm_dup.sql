delete from students
where rowid not in (
    select  min(rowid)
    from    students
    group by id
)