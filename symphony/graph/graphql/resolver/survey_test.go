// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint: goconst
package resolver

import (
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestAddRemoveSurvey(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, sr, wfr, cellr := r.Mutation(), r.Query(), r.Survey(), r.SurveyWiFiScan(), r.SurveyCellScan()

	locationType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "location_type_name_1",
	})
	require.NoError(t, err)
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name_1",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	ctyp := models.SurveyQuestionTypeCellular
	wftyp := models.SurveyQuestionTypeWifi

	bsID := "BSID-1"
	timestamp := 1564523995
	cellData := models.SurveyCellScanData{
		NetworkType:    models.CellularNetworkTypeGsm,
		SignalStrength: 1,
		BaseStationID:  &bsID,
		Timestamp:      &timestamp,
	}
	resp1 := models.SurveyQuestionResponse{
		FormIndex:      0,
		QuestionIndex:  0,
		QuestionFormat: &ctyp,
		QuestionText:   "give cellular data?",
		CellData:       []*models.SurveyCellScanData{&cellData},
	}
	wifiData := models.SurveyWiFiScanData{
		Timestamp: timestamp,
		Frequency: 1,
		Channel:   2,
		Bssid:     "bssid 2",
		Strength:  3,
	}
	resp2 := models.SurveyQuestionResponse{
		FormIndex:      1,
		QuestionIndex:  1,
		QuestionFormat: &wftyp,
		QuestionText:   "give wifi data?",
		WifiData:       []*models.SurveyWiFiScanData{&wifiData},
	}
	basicSurveyData := models.SurveyCreateData{
		Name:                "survey_one",
		LocationID:          location.ID,
		CompletionTimestamp: 100000,
		SurveyResponses:     []*models.SurveyQuestionResponse{&resp1, &resp2},
	}
	surveyID, err := mr.CreateSurvey(ctx, basicSurveyData)
	require.NoError(t, err)

	file, err := mr.AddImage(ctx, models.AddImageInput{
		EntityType:  models.ImageEntitySiteSurvey,
		EntityID:    *surveyID,
		ImgKey:      "1234",
		FileName:    "site_survey.",
		FileSize:    50000,
		Modified:    time.Now(),
		ContentType: "application/zip",
	})
	require.NoError(t, err)

	surveys, err := qr.Surveys(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, surveys, 1, "Verifying 'Surveys' return value")
	fetchedSurvey := surveys[0]

	require.Equal(t, *surveyID, fetchedSurvey.ID, "Verifying saved survey vs fetched survey: ID")

	slID, err := sr.LocationID(ctx, fetchedSurvey)
	require.NoError(t, err)
	require.Equal(t, location.ID, slID, "VVerifying saved location vs fetched location: locationType")

	require.Equal(t, "survey_one", fetchedSurvey.Name, "Verifying saved survey vs fetched survey: Name")

	ts, err := sr.CompletionTimestamp(ctx, fetchedSurvey)
	require.NoError(t, err)
	require.Equal(t, 100000, ts, "Verifying saved survey vs fetched survey: CompletionTimestamp")

	responses, err := sr.SurveyResponses(ctx, fetchedSurvey)
	require.NoError(t, err)
	require.Len(t, responses, 2, "Verifying nu,ber of responses")

	cellQ := fetchedSurvey.QueryQuestions().Where(surveyquestion.QuestionFormat(models.SurveyQuestionTypeCellular.String())).OnlyX(ctx)
	wfQ := fetchedSurvey.QueryQuestions().Where(surveyquestion.QuestionFormat(models.SurveyQuestionTypeWifi.String())).OnlyX(ctx)

	nt, err := cellr.NetworkType(ctx, cellQ.QueryCellScan().OnlyX(ctx))
	require.NoError(t, err)
	require.Equal(t, 0, cellQ.FormIndex, "Verifying saved survey vs fetched survey: FormIndex")
	require.Equal(t, 0, cellQ.QuestionIndex, "Verifying saved survey vs fetched survey: QuestionIndex")
	require.Equal(t, models.CellularNetworkTypeGsm, nt, "Verifying fetched survey cell scan's network type")

	cellTmstmp, err := cellr.Timestamp(ctx, cellQ.QueryCellScan().OnlyX(ctx))
	require.NoError(t, err)
	require.NotNil(t, cellTmstmp, "Verifying fetched survey cell scan's timestamp is not nil")
	require.Equal(t, timestamp, *cellTmstmp, "Verifying fetched survey cell scan's timestamp is correct")

	wfTmstmp, err := wfr.Timestamp(ctx, wfQ.QueryWifiScan().OnlyX(ctx))
	require.NoError(t, err)
	require.Equal(t, 1, wfQ.FormIndex, "Verifying saved survey vs fetched survey: FormIndex")
	require.Equal(t, 1, wfQ.QuestionIndex, "Verifying saved survey vs fetched survey: QuestionIndex")
	require.Equal(t, 3, wfQ.QueryWifiScan().OnlyX(ctx).Strength, "Verifying fetched survey wifi scan's Strength")
	require.Equal(t, timestamp, wfTmstmp, "Verifying fetched survey wifi scan's timestamp")

	fetchedFile, err := fetchedSurvey.QuerySourceFile().Only(ctx)
	require.NoError(t, err)
	require.Equal(t, file.ID, fetchedFile.ID, "Source file of site survey not fetched well")

	_, err = mr.RemoveSiteSurvey(ctx, fetchedSurvey.ID)
	require.NoError(t, err)
	surveys, err = qr.Surveys(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, surveys, 0, "Verifying 'Surveys' return value")

	surveyQuestionsExist, err := r.client.SurveyQuestion.Query().Exist(ctx)
	require.NoError(t, err)
	require.False(t, surveyQuestionsExist, "Survey questions still exit")

	surveyCellScansExist, err := r.client.SurveyCellScan.Query().Exist(ctx)
	require.NoError(t, err)
	require.False(t, surveyCellScansExist, "Survey cell scans still exit")

	surveyWifiScansExist, err := r.client.SurveyWiFiScan.Query().Exist(ctx)
	require.NoError(t, err)
	require.False(t, surveyWifiScansExist, "Survey wifi scans still exit")
}
