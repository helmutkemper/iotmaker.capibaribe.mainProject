golbalOnInit.push(function(){
  $("#containerConfigurationImageName").kendoComboBox({
    dataSource: new kendo.data.DataSource({
      serverFiltering: false,
      transport: {
        read: function(options) {
          $.ajax({
            type: "GET",
            url: "http://panel.localhost:8888/imageList",
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            //data: JSON.stringify({key: "value"}),
            success: function(data) {
              console.log("success");
              var newData = [];
              for( var i in data.Objects ){
                i = parseInt(i);
                newData.push({
                  id: data.Objects[i].Id,
                  name: data.Objects[i].RepoTags[0],
                  created: new Date( data.Objects[i].Created * 1000 ),
                  labels: data.Objects[i].Labels,
                  repoDigests: data.Objects[i].RepoDigests,
                  repoTags: data.Objects[i].RepoTags,
                  size: data.Objects[i].Size,
                  virtualSize: data.Objects[i].VirtualSize,
                });
              }
              console.log("new data:",newData);
              options.success(newData);
            }
          });
        }
      }
    }),
    change: function(e) {
      console.log("this.init.context.value:", this.element.context.value);
      // fixme: pegar a lista de configurações  
    },
    dataTextField: "name",
    dataValueField: "id",
    width: 624
  });
  containerConfigurationImageNameRef = $("#containerConfigurationImageName").data("kendoComboBox");
});