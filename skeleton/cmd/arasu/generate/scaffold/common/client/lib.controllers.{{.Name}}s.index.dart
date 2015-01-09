part of {{.App.Name}}.controllers.{{.Name}}s;
@Component(selector: '{{.Name}}s-index',  
  templateUrl: 'packages/{{.App.Name}}/views/{{.Name}}s/index.html')

class {{.Cname}}sIndex extends Cntr { 
  List<{{.Cname}}> records;

  {{.Cname}}sIndex(RouteProvider routeProvider,
    Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready..then((_){
      store.once();
      records = store.caches;
    })
    ..catchError(Logger.root.shout);
  }
  
  void destroy(var record){
    store.remove(record);
  }
}
