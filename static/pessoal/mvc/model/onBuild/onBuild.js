golbalOnInit.push(function(){
  containerHostOnBuildAddNewOnBuildDialogWindowRef = $("#containerHostOnBuildAddNewOnBuild");

  $("#containerHostOnBuild").kendoMultiSelect({
    filter: "contains",
    placeholder: "",
    itemTemplate: kendo.template($("#containerHostOnBuildTemplate").html()),
    noDataTemplate: kendo.template($("#containerHostOnBuildNoDataTemplate").html()),
    footerTemplate: kendo.template($("#containerHostOnBuildFooterTemplate").html()),
    dataTextField: "OnBuild",
    dataValueField: "Id",
    dataSource: new kendo.data.DataSource({
      schema: {
        model: {
          id: "Id",
          fields: {
            Id: { type: "number" },
            OnBuild: { type: "string" }
          }
        }
      }
    })
  });
  containerHostOnBuildRef = $("#containerHostOnBuild").data("kendoMultiSelect");
});

var containerConfigurationOnBuildNameRef;
var containerHostOnBuildRef;
var containerHostOnBuildAddNewOnBuildDialogWindowRef;

var containerHostOnBuildItemsCount = -1;
var containerHostOnBuildItemsIdToAdd = -1;

function containerHostOnBuildAddNewOnBuild(){
  containerHostOnBuildAddNewOnBuildDialogWindowRef.kendoDialog({
    title: "On build metadata from container",
    content: kendo.template($("#containerCreateTemplateOnBuildAddNewOnBuild").html()),
    visible: false,
    modal: true,
    close: function(){

    },
    open: function(){
      containerConfigurationOnBuildNameRef = $("#OnBuildName");

      setTimeout( function(){ containerConfigurationOnBuildNameRef.focus(); }, 1000)
    },
    actions: [
      {
        text: "Close"
      },
      {
        text: "Add",
        action: function(e){
          containerHostOnBuildAddNewOnBuildFunction();
          return false;
        }
      },
      {
        text: "Add and close",
        action: function(e){
          containerHostOnBuildAddNewOnBuildFunction();
        },
        primary: true
      }
    ]
  });
  containerHostOnBuildAddNewOnBuildDialogWindowRef.data("kendoDialog").open();
}

function containerHostOnBuildAddNewOnBuildFunction(){


  let dataSource = containerHostOnBuildRef.dataSource;
  containerHostOnBuildItemsCount += 1;
  containerHostOnBuildItemsIdToAdd = containerHostOnBuildItemsCount;

  let dataInDataSource = dataSource.data();
  let pass = true;

  for (var i in dataInDataSource) {
    if( isNaN( i ) ){
      break
    }
    i = parseInt(i);
    if (dataInDataSource[i].OnBuild === $("#OnBuildName").val()) {
      containerHostOnBuildItemsIdToAdd = dataInDataSource[i].Id;
      pass = false;
      break;
    }
  }

  if( pass === true ) {
    dataSource.add({
      Id: containerHostOnBuildItemsCount,
      OnBuild: $("#OnBuildName").val()
    });
  }

  dataSource.one("requestEnd", function(args) {
    if (args.type !== "create") {
      return;
    }

    dataSource.one("sync", function() {
      containerHostOnBuildRef.value(containerHostOnBuildRef.value().concat([containerHostOnBuildItemsIdToAdd]));
    });

    $("#OnBuildName").val("");
  });

  dataSource.sync();
}
