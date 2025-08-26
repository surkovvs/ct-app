package ctapp

import (
	"context"
	"reflect"

	"github.com/surkovvs/ct-app/component"
)

func (a *app) Start(ctx context.Context) {
	a.exec(ctx)
	<-a.shutdown.shutdownDone
}

func (a *app) AddModule(moduleName string, module any) {
	groupName := nameficator.getNextGroupName()
	comp := component.DefineComponent(component.Define{
		GroupName:  groupName,
		CompName:   moduleName,
		Component:  module,
		ExecWG:     a.execution.wg,
		ShutdownWG: a.shutdown.wg,
		ErrChan:    a.execution.errFlow,
	})
	if !comp.IsValid() {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, groupName,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, "module does not implement valid methods")
		return
	}
	if err := a.storage.AddComponent(groupName, moduleName, comp); err != nil {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, groupName,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, err)
	}
}

func (a *app) AddModuleToGroup(groupName, moduleName string, module any) {
	comp := component.DefineComponent(component.Define{
		GroupName:  groupName,
		CompName:   moduleName,
		Component:  module,
		ExecWG:     a.execution.wg,
		ShutdownWG: a.shutdown.wg,
		ErrChan:    a.execution.errFlow,
	})
	if !comp.IsValid() {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, groupName,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, "module does not implement valid methods")
		return
	}
	if err := a.storage.AddComponent(groupName, moduleName, comp); err != nil {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, groupName,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, err)
	}
}

func (a *app) AddBackgroundModule(moduleName string, module any) {
	comp := component.DefineComponent(component.Define{
		GroupName:  BackgroundGroup,
		CompName:   moduleName,
		Component:  module,
		ExecWG:     a.execution.wg,
		ShutdownWG: a.shutdown.wg,
		ErrChan:    a.execution.errFlow,
	})
	if !comp.IsValid() {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, BackgroundGroup,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, "module does not implement valid methods")
		return
	}
	if err := a.storage.AddComponent(BackgroundGroup, moduleName, comp); err != nil {
		a.logger.Error(`module addition`,
			"application", a.name,
			`group`, BackgroundGroup,
			`module`, moduleName,
			`unapplyed`, reflect.ValueOf(module).Type().Name(),
			`error`, err)
	}
}
