package main

import "github.com/backend-common/common"

func saveDummyData() {
	var size int64
	if err := db.Model(&Question{}).Count(&size).Error; err != nil {
		log.Error("", err)
		return
	} else if size >= 12 {
		// Already populated
		return
	}

	questionA := &Question{
		Difficulty: common.DIFFICULTY_EASY,
		Name:       "Two Sum",
		Description: `Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
		You may assume that each input would have exactly one solution, and you may not use the same element twice.
		You can return the answer in any order.`,
	}
	questionB := &Question{
		Difficulty: common.DIFFICULTY_EASY,
		Name:       "Palindrome Number",
		Description: `Given an integer x, return true if x is palindrome integer.
		An integer is a palindrome when it reads the same backward as forward.For example, 121 is a palindrome while 123 is not.`,
	}
	questionC := &Question{
		Difficulty:  common.DIFFICULTY_EASY,
		Name:        "Roman to Integer",
		Description: `Given a roman numeral, convert it to an integer.`,
	}
	questionD := &Question{
		Difficulty:  common.DIFFICULTY_MEDIUM,
		Name:        "Add Two Numbers",
		Description: `You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.`,
	}
	questionE := &Question{
		Difficulty:  common.DIFFICULTY_HARD,
		Name:        "Median of Two Sorted Arrays",
		Description: `Median of Two Sorted Arrays', 'Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays.`,
	}
	questionF := &Question{
		Difficulty: common.DIFFICULTY_EASY,
		Name:       "Two Sum",
		Description: `Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
		You may assume that each input would have exactly one solution, and you may not use the same element twice.
		You can return the answer in any order.`,
	}
	questionG := &Question{
		Difficulty: common.DIFFICULTY_EASY,
		Name:       "Longest Common Prefix",
		Description: `Write a function to find the longest common prefix string amongst an array of strings.
		If there is no common prefix, return an empty string "".`,
	}
	questionH := &Question{
		Difficulty:  common.DIFFICULTY_EASY,
		Name:        "Search Insert Position",
		Description: `Given a sorted array of distinct integers and a target value, return the index if the target is found. If not, return the index where it would be if it were inserted in order.`,
	}
	questionI := &Question{
		Difficulty:  common.DIFFICULTY_EASY,
		Name:        "Valid Parentheses",
		Description: `Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.`,
	}
	questionJ := &Question{
		Difficulty:  common.DIFFICULTY_EASY,
		Name:        "Remove Element",
		Description: `Given an integer array nums and an integer val, remove all occurrences of val in nums in-place. The relative order of the elements may be changed.`,
	}
	questionK := &Question{
		Difficulty:  common.DIFFICULTY_EASY,
		Name:        "Remove Duplicates from Sorted Array",
		Description: `Given an integer array nums sorted in non-decreasing order, remove the duplicates in-place such that each unique element appears only once. The relative order of the elements should be kept the same.`,
	}
	questionL := &Question{
		Difficulty: common.DIFFICULTY_EASY,
		Name:       "Merge Two Sorted Lists",
		Description: `You are given the heads of two sorted linked lists list1 and list2.
		Merge the two lists in a one sorted list. The list should be made by splicing together the nodes of the first two lists.
		Return the head of the merged linked list.`,
	}

	db.Create(questionA)
	db.Create(questionB)
	db.Create(questionC)
	db.Create(questionD)
	db.Create(questionE)
	db.Create(questionF)
	db.Create(questionG)
	db.Create(questionH)
	db.Create(questionI)
	db.Create(questionJ)
	db.Create(questionK)
	db.Create(questionL)
}
