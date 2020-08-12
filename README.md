<!--
 * @Author: jinde.zgm
 * @Date: 2020-08-12 21:59:48
 * @Descripttion: 
-->
# Introduce
Concurrent is a package designed for concurrency.
## Map
Map extend Update interface based on sync.Map, equivalent to CAS(Compare And Swap).  
Update keeps calling 'tryUpdate()' to update key 'key' retrying the update until success if there is conflict. Note that value passed to tryUpdate may change across invocations of tryUpdate() if other writers are simultaneously updating it, so tryUpdate() needs to take into account the current contents of the value when deciding how the update object should look. If the key doesn't exist, it will pass nil to tryUpdate.