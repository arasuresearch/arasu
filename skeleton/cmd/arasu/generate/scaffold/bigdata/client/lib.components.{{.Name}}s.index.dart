part of components.{{.Name}}s;
@Component(selector: 'index',  templateUrl: 'packages/{{.App.Name}}/components/{{.Name}}s/index.html', publishAs: 'C')
class Index extends Cntr { 
  List<{{.Cname}}> records;

  Index(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
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


