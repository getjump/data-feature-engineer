package features
//
//import (
//	"testing"
//)
//
//func TestFeatureEngineer_AppendFeature(t *testing.T) {
//	f := &MinFeature{}
//	featureEngineer := FeatureEngineer{}
//
//	featureEngineer.AppendFeature(f)
//
//	if len(featureEngineer.Features) != 1 && featureEngineer.Features[0] != f {
//		t.Errorf("FeatureEngineer.AppendFeature: Feature was not appended")
//	}
//}
//
//type MockFeature struct {
//	CalledTimes uint64
//	BasicFeature
//}
//func (m *MockFeature) Update(TimeCurrent uint64, data []*InputData) {
//	m.CalledTimes++
//}
//
//func TestFeatureEngineer_Update(t *testing.T) {
//	featureEngineer := FeatureEngineer{}
//	f := &MockFeature{CalledTimes: 0}
//
//	featureEngineer.AppendFeature(f)
//	featureEngineer.Update(nil)
//
//	if f.CalledTimes != 1 {
//		t.Errorf("FeatureEngineer.Update: Havent called feature update")
//	}
//}