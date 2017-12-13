{{template "base/base.html" .}}
{{define "body"}}
  <div class="right-content-container">
    <div class="header">
      <ol class="breadcrumb">
        <li>
          <a href="/">Home</a>
        </li>
        <li class="active">Settings</li>
      </ol>
    </div>
    <div class="row">
      <div class="col-md-12">
        <h1>About</h1>
      </div>
    </div>
    <div class="content-block white-bg">
      <div class="row" style="height:25px;">
        <div class="pull-right">
          <a href="https://github.com/snagles/docker-registry-manager"><img src="https://img.shields.io/github/stars/snagles/docker-registry-manager.svg?style=social&amp;label=Star" alt="GitHub stars"/></a>
          <a href="https://github.com/snagles/docker-registry-manager/issues"><img src="https://img.shields.io/github/issues-raw/snagles/docker-registry-manager.svg" alt="GitHub issues"/></a>
          <a href="https://raw.githubusercontent.com/snagles/docker-registry-manager/master/LICENSE"><img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license"></a>
        </div>
      </div>
      <div class="row">
        <table class="table table-striped table-bordered">
          <thead>
            <tr>
              <th>Service</th>
              <th>Master</th>
              <th>Develop</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Status</td>
              <td>
                <a href="https://github.com/snagles/docker-registry-manager/tree/master"><img src="https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master" alt="Build Status"/></a>
              </td>
              <td>
                <a href="https://github.com/snagles/docker-registry-manager/tree/develop"><img src="https://travis-ci.org/snagles/docker-registry-manager.svg?branch=develop" alt="Build Status"/></a>
              </td>
            </tr>
            <tr>
              <td>Coverage</td>
              <td>
                <a href="https://codecov.io/gh/snagles/docker-registry-manager">
                  <img src="https://codecov.io/gh/snagles/docker-registry-manager/branch/master/graph/badge.svg" alt="Coverage Status"></a>
              </td>
              <td>
                <a href="https://codecov.io/gh/snagles/docker-registry-manager">
                  <img src="https://codecov.io/gh/snagles/docker-registry-manager/branch/develop/graph/badge.svg" alt="Coverage Status"></a>
              </td>
            </tr>
            <tr>
              <td>Documentation</td>
              <td>
                <a href="https://godoc.org/github.com/snagles/docker-registry-manager">
                  <img src="https://godoc.org/github.com/snagles/docker-registry-manager?status.svg" alt="GoDoc"></a>
              </td>
              <td>
                <a href="https://godoc.org/github.com/snagles/docker-registry-manager">
                  <img src="https://godoc.org/github.com/snagles/docker-registry-manager?status.svg" alt="GoDoc"></a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
{{ end }}
