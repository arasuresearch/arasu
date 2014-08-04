part of components.{{.Name}}s;
@Component(selector: 'show',  templateUrl: 'packages/{{.App.Name}}/components/{{.Name}}s/show.html', publishAs: 'C')
class Show extends Cntr{
  {{.Cname}} record;

  Show(RouteProvider routeProvider,Db db,NgRoutingHelper locationService) :
    super({{.Cname}},routeProvider,db,locationService) {
    ready.then((_){
      record = store.find(Id);//store.find(int.parse(Id, onError: (_) => 0));
    }).catchError((e)=>Logger.root.shout);
  }
  void destroy(var record){
    store.remove(record);
    locationService.router.gotoUrl('/');
  }
}
