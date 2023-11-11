package models

import (
	"errors"
	"reflect"
	"sync"
)

var modelManager = NewManager(nil)

type ManagerConfig struct {
	DatabaseName string
}

type Manager struct {
	cfg       *ManagerConfig
	models    map[reflect.Type]*Declaration
	reflector *Reflector
	mutex     sync.RWMutex
}

func NewManager(cfg *ManagerConfig) *Manager {
	toReturn := &Manager{
		cfg:       cfg,
		models:    make(map[reflect.Type]*Declaration),
		reflector: &Reflector{},
	}

	if cfg != nil {
		toReturn.reflector.databaseName = cfg.DatabaseName
	}

	return toReturn
}

func (m *Manager) GetConfig() *ManagerConfig {
	return m.cfg
}

func (m *Manager) SetConfig(cfg *ManagerConfig) {
	m.cfg = cfg
	m.reflector.databaseName = cfg.DatabaseName
}

func (m *Manager) RegisterModel(model any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	modelType := reflect.TypeOf(model).Elem()

	if _, exists := m.models[modelType]; !exists {
		declaration, err := m.reflector.ReflectModel(modelType)
		if err != nil {
			return err
		}

		m.models[modelType] = declaration
	}

	return nil
}

func (m *Manager) GetDeclaration(model any) (*Declaration, error) {
	modelType := reflect.TypeOf(model).Elem()
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	declaration, exists := m.models[modelType]
	if !exists {
		return nil, errors.New("Model not registered")
	}

	return declaration, nil
}

func M() *Manager {
	return modelManager
}
