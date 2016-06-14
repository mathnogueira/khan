// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"testing"

	"github.com/Pallinder/go-randomdata"
	. "github.com/franela/goblin"
)

func TestMembershipModel(t *testing.T) {
	g := Goblin(t)
	testDb, err := GetTestDB()
	g.Assert(err == nil).IsTrue()

	g.Describe("Membership Model", func() {
		g.It("Should create a new Membership", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			membership := &Membership{
				GameID:      "test",
				ClanID:      clan.ID,
				PlayerID:    player.ID,
				RequestorID: player.ID,
				Level:       1,
				Approved:    false,
				Denied:      false,
			}
			err = testDb.Insert(membership)
			g.Assert(err == nil).IsTrue()
			g.Assert(membership.ID != 0).IsTrue()

			dbMembership, err := GetMembershipByID(testDb, membership.ID)
			g.Assert(err == nil).IsTrue()

			g.Assert(dbMembership.GameID).Equal(membership.GameID)
			g.Assert(dbMembership.PlayerID).Equal(membership.PlayerID)
			g.Assert(dbMembership.ClanID).Equal(membership.ClanID)
		})

		g.It("Should update a Membership", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			membership := MembershipFactory.MustCreateWithOption(map[string]interface{}{
				"PlayerID":    player.ID,
				"ClanID":      clan.ID,
				"RequestorID": clan.OwnerID,
			}).(*Membership)
			err = testDb.Insert(membership)
			g.Assert(err == nil).IsTrue()
			dt := membership.UpdatedAt

			membership.Approved = true
			count, err := testDb.Update(membership)
			g.Assert(err == nil).IsTrue()
			g.Assert(int(count)).Equal(1)
			g.Assert(membership.UpdatedAt > dt).IsTrue()
		})

		g.It("Should get existing Membership", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			membership := MembershipFactory.MustCreateWithOption(map[string]interface{}{
				"PlayerID":    player.ID,
				"ClanID":      clan.ID,
				"RequestorID": player.ID,
			}).(*Membership)
			err = testDb.Insert(membership)
			g.Assert(err == nil).IsTrue()

			dbMembership, err := GetMembershipByID(testDb, membership.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbMembership.ID).Equal(membership.ID)
		})

		g.It("Should not get non-existing Membership", func() {
			_, err = GetMembershipByID(testDb, -1)
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("Membership was not found with id: -1")
		})

		g.It("Should get an existing Membership by the player public ID", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  player.GameID,
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			membership := MembershipFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":      player.GameID,
				"PlayerID":    player.ID,
				"ClanID":      clan.ID,
				"RequestorID": player.ID,
			}).(*Membership)
			err = testDb.Insert(membership)
			g.Assert(err == nil).IsTrue()

			dbMembership, err := GetMembershipByPlayerPublicID(testDb, player.GameID, player.PublicID)
			g.Assert(err == nil).IsTrue()
			g.Assert(dbMembership.ID).Equal(membership.ID)
			g.Assert(dbMembership.PlayerID).Equal(player.ID)
		})

		g.It("Should not get non-existing Membership by the player public ID", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  player.GameID,
				"OwnerID": player.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			dbMembership, err := GetMembershipByPlayerPublicID(testDb, player.GameID, player.PublicID)
			g.Assert(err != nil).IsTrue()
			g.Assert(dbMembership == nil).IsTrue()
		})

		g.It("Should not get non-existing Membership", func() {
			_, err = GetMembershipByID(testDb, -1)
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("Membership was not found with id: -1")
		})

		g.Describe("Should create a new Membership with CreateMembership", func() {

			g.It("If requestor is the player", func() {
				player := PlayerFactory.MustCreate().(*Player)
				err := testDb.Insert(player)
				g.Assert(err == nil).IsTrue()

				owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
					"GameID": player.GameID,
				}).(*Player)
				err = testDb.Insert(owner)
				g.Assert(err == nil).IsTrue()

				clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
					"GameID":  owner.GameID,
					"OwnerID": owner.ID,
				}).(*Clan)
				err = testDb.Insert(clan)
				g.Assert(err == nil).IsTrue()

				membership, err := CreateMembership(
					testDb,
					player.GameID,
					1,
					player.PublicID,
					clan.PublicID,
					player.PublicID,
				)

				g.Assert(err == nil).IsTrue()
				g.Assert(membership.ID != 0).IsTrue()

				dbMembership, err := GetMembershipByID(testDb, membership.ID)
				g.Assert(err == nil).IsTrue()

				g.Assert(dbMembership.GameID).Equal(membership.GameID)
				g.Assert(dbMembership.PlayerID).Equal(player.ID)
				g.Assert(dbMembership.RequestorID).Equal(player.ID)
				g.Assert(dbMembership.ClanID).Equal(clan.ID)
				g.Assert(dbMembership.Approved).Equal(false)
				g.Assert(dbMembership.Denied).Equal(false)
			})

			g.It("If requestor is the clan owner", func() {
				player := PlayerFactory.MustCreate().(*Player)
				err := testDb.Insert(player)
				g.Assert(err == nil).IsTrue()

				owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
					"GameID": player.GameID,
				}).(*Player)
				err = testDb.Insert(owner)
				g.Assert(err == nil).IsTrue()

				clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
					"GameID":  owner.GameID,
					"OwnerID": owner.ID,
				}).(*Clan)
				err = testDb.Insert(clan)
				g.Assert(err == nil).IsTrue()

				membership, err := CreateMembership(
					testDb,
					player.GameID,
					1,
					player.PublicID,
					clan.PublicID,
					owner.PublicID,
				)

				g.Assert(err == nil).IsTrue()
				g.Assert(membership.ID != 0).IsTrue()

				dbMembership, err := GetMembershipByID(testDb, membership.ID)
				g.Assert(err == nil).IsTrue()

				g.Assert(dbMembership.GameID).Equal(membership.GameID)
				g.Assert(dbMembership.PlayerID).Equal(player.ID)
				g.Assert(dbMembership.RequestorID).Equal(owner.ID)
				g.Assert(dbMembership.ClanID).Equal(clan.ID)
				g.Assert(dbMembership.Approved).Equal(false)
				g.Assert(dbMembership.Denied).Equal(false)
			})

			g.It("If requestor is a member of the clan with level greater than the min level", func() {
				player := PlayerFactory.MustCreate().(*Player)
				err := testDb.Insert(player)
				g.Assert(err == nil).IsTrue()

				owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
					"GameID": player.GameID,
				}).(*Player)
				err = testDb.Insert(owner)
				g.Assert(err == nil).IsTrue()

				requestor := PlayerFactory.MustCreateWithOption(map[string]interface{}{
					"GameID": player.GameID,
				}).(*Player)
				err = testDb.Insert(requestor)
				g.Assert(err == nil).IsTrue()

				clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
					"GameID":  owner.GameID,
					"OwnerID": owner.ID,
				}).(*Clan)
				err = testDb.Insert(clan)
				g.Assert(err == nil).IsTrue()

				requestorMembership := &Membership{
					GameID:      player.GameID,
					ClanID:      clan.ID,
					PlayerID:    requestor.ID,
					RequestorID: owner.ID,
					Level:       5,
					Approved:    true,
					Denied:      false,
				}
				err = testDb.Insert(requestorMembership)

				membership, err := CreateMembership(
					testDb,
					player.GameID,
					1,
					player.PublicID,
					clan.PublicID,
					requestor.PublicID,
				)

				g.Assert(err == nil).IsTrue()
				g.Assert(membership.ID != 0).IsTrue()

				dbMembership, err := GetMembershipByID(testDb, membership.ID)
				g.Assert(err == nil).IsTrue()

				g.Assert(dbMembership.GameID).Equal(membership.GameID)
				g.Assert(dbMembership.PlayerID).Equal(player.ID)
				g.Assert(dbMembership.RequestorID).Equal(requestor.ID)
				g.Assert(dbMembership.ClanID).Equal(clan.ID)
				g.Assert(dbMembership.Approved).Equal(false)
				g.Assert(dbMembership.Denied).Equal(false)
			})
		})

		g.It("Should not create a new Membership with CreateMembership if unexistent player", func() {
			owner := PlayerFactory.MustCreate().(*Player)
			err = testDb.Insert(owner)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  owner.GameID,
				"OwnerID": owner.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			playerPublicID := randomdata.FullName(randomdata.RandomGender)
			_, err = CreateMembership(
				testDb,
				owner.GameID,
				1,
				playerPublicID,
				clan.PublicID,
				owner.PublicID,
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Player was not found with id: %s", playerPublicID))
		})

		g.It("Should not create a new Membership with CreateMembership if unexistent clan", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			clanPublicID := randomdata.FullName(randomdata.RandomGender)
			_, err = CreateMembership(
				testDb,
				player.GameID,
				1,
				player.PublicID,
				clanPublicID,
				player.PublicID,
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Clan was not found with id: %s", clanPublicID))

		})

		g.It("Should not create a new Membership with CreateMembership if unexistent requestor", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
				"GameID": player.GameID,
			}).(*Player)
			err = testDb.Insert(owner)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  owner.GameID,
				"OwnerID": owner.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			requestorPublicID := randomdata.FullName(randomdata.RandomGender)
			_, err = CreateMembership(
				testDb,
				player.GameID,
				1,
				player.PublicID,
				clan.PublicID,
				requestorPublicID,
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Player %s cannot create membership for clan %s", requestorPublicID, clan.PublicID))
		})

		g.It("Should not create a new Membership with CreateMembership if requestor's level is less than min level", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
				"GameID": player.GameID,
			}).(*Player)
			err = testDb.Insert(owner)
			g.Assert(err == nil).IsTrue()

			requestor := PlayerFactory.MustCreateWithOption(map[string]interface{}{
				"GameID": player.GameID,
			}).(*Player)
			err = testDb.Insert(requestor)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  owner.GameID,
				"OwnerID": owner.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			requestorMembership := &Membership{
				GameID:      player.GameID,
				ClanID:      clan.ID,
				PlayerID:    requestor.ID,
				RequestorID: owner.ID,
				Level:       0,
				Approved:    true,
				Denied:      false,
			}
			err = testDb.Insert(requestorMembership)

			_, err = CreateMembership(
				testDb,
				player.GameID,
				1,
				player.PublicID,
				clan.PublicID,
				requestor.PublicID,
			)

			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal(fmt.Sprintf("Player %s cannot create membership for clan %s", requestor.PublicID, clan.PublicID))
		})

		g.It("Should fail if membership already exists", func() {
			player := PlayerFactory.MustCreate().(*Player)
			err := testDb.Insert(player)
			g.Assert(err == nil).IsTrue()

			owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
				"GameID": player.GameID,
			}).(*Player)
			err = testDb.Insert(owner)
			g.Assert(err == nil).IsTrue()

			clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
				"GameID":  owner.GameID,
				"OwnerID": owner.ID,
			}).(*Clan)
			err = testDb.Insert(clan)
			g.Assert(err == nil).IsTrue()

			previousMembership := &Membership{
				GameID:      player.GameID,
				ClanID:      clan.ID,
				PlayerID:    player.ID,
				RequestorID: player.ID,
				Level:       0,
				Approved:    true,
				Denied:      false,
			}
			err = testDb.Insert(previousMembership)

			membership, err := CreateMembership(
				testDb,
				player.GameID,
				1,
				player.PublicID,
				clan.PublicID,
				player.PublicID,
			)

			g.Assert(membership == nil).IsTrue()
			g.Assert(err != nil).IsTrue()
			g.Assert(err.Error()).Equal("pq: duplicate key value violates unique constraint \"playerid_clanid\"")
		})
	})
}
