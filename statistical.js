$(document).ready(function () {
//https://www.cnblogs.com/adolfmc/p/7698364.html
    $.get(
        "http://localhost:9909/dig",
        {
            "time":gettime(),
            "url":geturl(),
            "referer":getrefer(),
            "ua":getuser_agent(),
        }
    );
})