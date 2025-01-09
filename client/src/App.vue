<template>
  <div id="luckysheet" style="margin:0px;padding:0px;position:absolute;width:100%;height:100%;left: 0px;top: 0px;"></div>
</template>

<script setup lang="ts" >
import {onMounted} from "vue";



const authToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Njg2LCJwaWQiOjEsImRlcHRJZCI6MTQ4LCJyb2xlSWQiOjIzNywicm9sZUtleSI6IjEyZWQiLCJyb2xlTGV2ZWwiOjAsImVtcGxveWVlSWQiOjgxMSwidXNlcm5hbWUiOiJrdXp1bWEiLCJyZWFsTmFtZSI6Imt1enVtYSIsImF2YXRhciI6Imh0dHBzOi8vdHotcmVzdW1lLm9zcy1hcC1zb3V0aGVhc3QtMS5hbGl5dW5jcy5jb20vdGlhbnpoaS9jb21tb24vZGVmYXVsdF9hdmF0YXIvMS5wbmclNDAzeCUyMCUyODIlMjkucG5nIiwiZW1haWwiOiIiLCJtb2JpbGUiOiIiLCJhcHAiOiJhZG1pbiIsImxvZ2luQXQiOiIyMDI0LTEwLTIzIDE3OjI5OjU4IiwidGVuYW50SWQiOjF9.VVDR7bu6O7S7T9ZVXGlwHcBfrDNFzLcwOrmwT8KIk0E'
setInterval(()=>{
},10000)

onMounted(()=>{
  luckysheet.create({
    container: 'luckysheet',
    allowUpdate: true,
    loadUrl: `http://172.16.19.7:2234/collaborate?authorization=${authToken}`,
    updateUrl: `ws://172.16.19.7:2234/collaborate?authorization=${authToken}`,
    hook: {
      updated: () => {
        // console.log(luckysheet.getAllSheets())
      }
      // updated: luckysheetUpdated,
      // cellRenderAfter: luckysheetCellRenderAfter,
    }
  })



  console.log(typeof luckysheet.ws)
})
</script>