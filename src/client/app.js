var app = angular.module('cd', ['ui.bootstrap']);

app.controller('MainCtrl', function($http, NodeSvc) {

  this.activeNode = function() {
    return NodeSvc.selectedNode;
  }

  this.nodes = function() {
    return NodeSvc.nodes;
  };

  this.click = function(node) {
    NodeSvc.selectNode(node);
  };
  
});

app.service('NodeSvc', function($http){
  this.nodes = null;
  this.activeNode = null;

  
  this.init = function() {
    var loadNodes = this.loadNodes.bind(this);
    $http.get('/cd/')
      .success(loadNodes);
  };

  this.loadNodes = function(data) {
    this.nodes = [];
    angular.forEach(data, function(item){
      this.nodes.push(new Node(item));
    }, this);
  };

  this.selectNode = function(node) {
    this.selectedNode = node;
    this.getChildNodes(node);
  };

  this.getChildNodes = function(node) {
    if (node.childUrl === null) {
      return;
    }
    if (node.loaded) {
      return;
    }
    $http.get(node.childUrl)
      .success(node.loadChildNodes.bind(node))
      .error(function(a, b, c, d){console.log(a, b, c, d)});
  };
  
  this.init();
});

/**
 * @param nodeData {!Object} item from response containing node data with
 * data, displayName, templateUrl, and childUrl properties.
 */
function Node(nodeData) {
  this.data = nodeData.data;
  this.displayName = nodeData.displayName;
  this.templateUrl = nodeData.templateUrl ? nodeData.templateUrl : null;
  this.childUrl = nodeData.childUrl ? nodeData.childUrl : null;
  
  this.open = false;
  this.loading = true;
  this.loaded = false;
  this.parentNode = null;
  
  this.internalChildNodes = [];
  this.leafChildNodes = [];

  this.loadChildNodes = function(data) {
      angular.forEach(data, function(item){
        this.addChildNode(new Node(item));
      }, this);
      this.loading=false;
      this.loaded=true;
  };

  this.addChildNode = function(childNode) {
    childNode.parentNode = this;
    if (childNode.childUrl) {
      this.internalChildNodes.push(childNode);
    } else {
      this.leafChildNodes.push(childNode);
    }
  };
  
  this.childNodesLoaded = function() { 
    return this.childNodes !== null;
  };

}