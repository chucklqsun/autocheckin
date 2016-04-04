# autocheckin  
for lazy guy finish checkin task in internet  
  
How to use:  
1. modify the api config in vendor.go, point to the task url you want to checkin  
2. add username/password into file "account"  
3. fetch cookie and store by Chrome or other browser(temporary solution,login module will be added later)  
4. save the cookie as ""vendorName.account.cookie" filename in your binary "autocheckin" dir  , eg: duokan.a1.cookie  -> duokan is vendor name, a1 is account
5. run ./autocheckin  
  
This project is just undergoing development. Any suggestion is welcome!   
  
"account" file:  
zimuzu.a1=i_am_username|i_am_password  

cookie file:  
zimuzu.a1.cookie  
duokan.a1.cookie  
duokan.a2.cookie  

ret looks like:

**Begin
Job ->  duokan : a1  
Job ->  duokan : a2  
Job ->  zimuzu : a1  
map[status:1 info:登录成功！ data:map[url_back:http://www.zimuzu.tv/]]  
exec  a1 http://www.zimuzu.tv/User/Login/ajaxLogin  success  
map[msg:今日已签到 result:500002]  
exec  a2 https://www.duokan.com/checkin/v0/checkin  success  
map[msg:今日已签到 result:500002]  
exec  a1 https://www.duokan.com/checkin/v0/checkin  success  
End  
**
